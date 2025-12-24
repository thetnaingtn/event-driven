package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"tickets/entity"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/jmoiron/sqlx"
)

type OpsBookingRepository struct {
	db *sqlx.DB
}

func NewOpsBookingRepository(db *sqlx.DB) *OpsBookingRepository {
	return &OpsBookingRepository{
		db: db,
	}
}

func (r *OpsBookingRepository) OnBookingMade(ctx context.Context, event *entity.BookingMade) error {
	readModel := entity.OpsBooking{
		BookingID:  event.BookingID,
		Tickets:    nil,
		LastUpdate: time.Now(),
		BookedAt:   event.Header.PublishedAt,
	}

	err := r.createReadModel(ctx, readModel)
	if err != nil {
		return err
	}

	return nil
}

func (r *OpsBookingRepository) OnTicketReceiptIssued(ctx context.Context, event *entity.TicketReceiptIssued) error {
	return r.updateReadModelByTicketID(ctx, event.TicketID, func(ticket entity.OpsTicket) (entity.OpsTicket, error) {
		ticket.ReceiptIssuedAt = event.IssuedAt
		ticket.ReceiptNumber = event.ReceiptNumber

		return ticket, nil
	})
}

func (r *OpsBookingRepository) OnTicketBookingConfirmed(ctx context.Context, event *entity.TicketBookingConfirmed) error {
	return r.updateReadModelByBookingID(
		ctx,
		event.BookingID,
		func(rm entity.OpsBooking) (entity.OpsBooking, error) {

			ticket, ok := rm.Tickets[event.TicketID]
			if !ok {
				// we are using zero-value of OpsTicket
				log.
					FromContext(ctx).
					With("ticket_id", event.TicketID).
					Debug("Creating ticket read model for ticket %s")
			}

			ticket.PriceAmount = event.Price.Amount
			ticket.PriceCurrency = event.Price.Currency
			ticket.CustomerEmail = event.CustomerEmail
			ticket.Status = "confirmed"

			for k := range rm.Tickets {
				fmt.Printf("ticket id: %d", k)
			}

			rm.Tickets[event.TicketID] = ticket

			return rm, nil
		},
	)
}

func (r *OpsBookingRepository) OnTicketRefunded(ctx context.Context, event *entity.TicketRefunded) error {
	return r.updateReadModelByTicketID(ctx, event.TicketID, func(ticket entity.OpsTicket) (entity.OpsTicket, error) {
		ticket.Status = "refunded"
		return ticket, nil
	})
}

func (r *OpsBookingRepository) OnTicketPrinted(ctx context.Context, event *entity.TicketPrinted) error {
	return r.updateReadModelByTicketID(ctx, event.TicketID, func(ticket entity.OpsTicket) (entity.OpsTicket, error) {
		ticket.PrintedAt = event.Header.PublishedAt
		ticket.PrintedFileName = event.FileName

		return ticket, nil
	})
}

func (r *OpsBookingRepository) createReadModel(
	ctx context.Context,
	booking entity.OpsBooking,
) (err error) {
	payload, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO 
		    read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO NOTHING; -- read model may be already updated by another event - we don't want to override
`, payload, booking.BookingID)

	if err != nil {
		return fmt.Errorf("could not create read model: %w", err)
	}

	return nil
}

func (r *OpsBookingRepository) updateReadModelByBookingID(
	ctx context.Context,
	bookingID string,
	updateFunc func(ticket entity.OpsBooking) (entity.OpsBooking, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByBookingID(ctx, bookingID, tx)
			if errors.Is(err, sql.ErrNoRows) {
				// events arrived out of order - it should spin until the read model is created
				return fmt.Errorf("read model for booking %s not exist yet", bookingID)
			} else if err != nil {
				return fmt.Errorf("could not find read model: %w", err)
			}

			updatedRm, err := updateFunc(rm)
			if err != nil {
				return err
			}

			return r.updateReadModel(ctx, tx, updatedRm)
		},
	)
}

func (r *OpsBookingRepository) updateReadModelByTicketID(
	ctx context.Context,
	ticketID string,
	updateFunc func(ticket entity.OpsTicket) (entity.OpsTicket, error),
) (err error) {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByTicketID(ctx, ticketID, tx)
			if errors.Is(err, sql.ErrNoRows) {
				// events arrived out of order - it should spin until the read model is created
				return fmt.Errorf("read model for ticket %s not exist yet", ticketID)
			} else if err != nil {
				return fmt.Errorf("could not find read model: %w", err)
			}

			ticket, _ := rm.Tickets[ticketID]

			updatedRm, err := updateFunc(ticket)
			if err != nil {
				return err
			}

			rm.Tickets[ticketID] = updatedRm

			return r.updateReadModel(ctx, tx, rm)
		},
	)
}

func (r *OpsBookingRepository) updateReadModel(
	ctx context.Context,
	tx *sqlx.Tx,
	rm entity.OpsBooking,
) error {
	rm.LastUpdate = time.Now()

	payload, err := json.Marshal(rm)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO 
			read_model_ops_bookings (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO UPDATE SET payload = excluded.payload;
		`, payload, rm.BookingID)
	if err != nil {
		return fmt.Errorf("could not update read model: %w", err)
	}

	return nil
}

func (r *OpsBookingRepository) findReadModelByTicketID(
	ctx context.Context,
	ticketID string,
	db dbExecutor,
) (entity.OpsBooking, error) {
	var payload []byte

	err := db.QueryRowContext(
		ctx,
		"SELECT payload FROM read_model_ops_bookings WHERE payload::jsonb -> 'tickets' ? $1",
		ticketID,
	).Scan(&payload)
	if err != nil {
		return entity.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r *OpsBookingRepository) findReadModelByBookingID(
	ctx context.Context,
	bookingID string,
	db dbExecutor,
) (entity.OpsBooking, error) {
	var payload []byte

	err := db.QueryRowContext(
		ctx,
		"SELECT payload FROM read_model_ops_bookings WHERE booking_id = $1",
		bookingID,
	).Scan(&payload)
	if err != nil {
		return entity.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r *OpsBookingRepository) unmarshalReadModelFromDB(payload []byte) (entity.OpsBooking, error) {
	var dbReadModel entity.OpsBooking
	if err := json.Unmarshal(payload, &dbReadModel); err != nil {
		return entity.OpsBooking{}, err
	}

	if dbReadModel.Tickets == nil {
		dbReadModel.Tickets = map[string]entity.OpsTicket{}
	}

	return dbReadModel, nil
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
