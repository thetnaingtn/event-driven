package db

import (
	"context"
	"database/sql"
	"fmt"
	"tickets/entity"
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *entity.Booking) error {
	stmt := `
		INSERT INTO bookings(booking_id, show_id, number_of_tickets, customer_email)
		VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email) ON CONFLICT DO NOTHING
	`

	return updateInTx(ctx, r.db, sql.LevelRepeatableRead, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, stmt, booking)
		if err != nil {
			return fmt.Errorf("failed to insert booking: %w", err)
		}

		outboxPublisher, err := outbox.NewPublisherForDB(ctx, tx)
		if err != nil {
			return err
		}

		if err := event.NewBus(outboxPublisher).Publish(ctx, &entity.BookingMade{
			Header:          entity.NewMessageHeader(),
			NumberOfTickets: booking.NumberOfTickets,
			BookingID:       booking.BookingID,
			CustomerEmail:   booking.CustomerEmail,
			ShowID:          booking.ShowID,
		}); err != nil {
			return fmt.Errorf("could not publish event: %w", err)
		}

		return nil
	})
}
