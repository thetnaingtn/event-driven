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

var (
	ErrNotEnoughSeats = fmt.Errorf("not enough seats available for booking")
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

	return updateInTx(ctx, r.db, sql.LevelSerializable, func(ctx context.Context, tx *sqlx.Tx) error {
		availableSeats := 0
		if err := tx.GetContext(ctx, &availableSeats, `
			SELECT
				number_of_tickets AS available_seats
			FROM
				shows
			WHERE
				show_id = $1
		`, booking.ShowID); err != nil {
			return fmt.Errorf("failed to check available seats: %w", err)
		}

		alreadyBookedSeats := 0
		if err := tx.GetContext(ctx, &alreadyBookedSeats, `
				SELECT
					COALESCE(SUM(number_of_tickets), 0) AS already_booked_seats
				FROM
					bookings
				WHERE
					show_id = $1
			`,
			booking.ShowID); err != nil {
			return fmt.Errorf("could not get already booked seats: %w", err)
		}

		if availableSeats-alreadyBookedSeats < booking.NumberOfTickets {
			return ErrNotEnoughSeats
		}

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
