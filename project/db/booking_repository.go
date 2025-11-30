package db

import (
	"context"
	"fmt"
	"tickets/entity"

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

	_, err := r.db.NamedExecContext(ctx, stmt, booking)
	if err != nil {
		return fmt.Errorf("failed to insert booking: %w", err)
	}

	return nil
}
