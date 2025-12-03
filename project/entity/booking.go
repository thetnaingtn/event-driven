package entity

import "github.com/google/uuid"

type CreateBookingRequest struct {
	BookingID       uuid.UUID `json:"booking_id"`
	EventID         uuid.UUID `json:"event_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
}

type Booking struct {
	BookingID       uuid.UUID `json:"booking_id" db:"booking_id"`
	ShowID          uuid.UUID `json:"show_id" db:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets" db:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email" db:"customer_email"`
}
