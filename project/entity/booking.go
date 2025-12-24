package entity

import (
	"time"

	"github.com/google/uuid"
)

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

type OpsBooking struct {
	BookingID uuid.UUID `json:"booking_id"` // from BookingMade event
	BookedAt  time.Time `json:"booked_at"`  // from BookingMade event

	Tickets map[string]OpsTicket `json:"tickets"` // Tickets added/updated by TicketBookingConfirmed, TicketRefunded, TicketPrinted, TicketReceiptIssued

	LastUpdate time.Time `json:"last_update"` // updated when read model is updated
}

type OpsTicket struct {
	PriceAmount   string `json:"price_amount"`   // from TicketBookingConfirmed event
	PriceCurrency string `json:"price_currency"` // from TicketBookingConfirmed event
	CustomerEmail string `json:"customer_email"` // from TicketBookingConfirmed event

	// Status should be set to "confirmed" or "refunded"
	Status string `json:"status"` // set to "confirmed" by TicketBookingConfirmed, "refunded" by TicketRefunded

	PrintedAt       time.Time `json:"printed_at"`        // from TicketPrinted event
	PrintedFileName string    `json:"printed_file_name"` // from TicketPrinted event

	ReceiptIssuedAt time.Time `json:"receipt_issued_at"` // from TicketReceiptIssued event
	ReceiptNumber   string    `json:"receipt_number"`    // from TicketReceiptIssued event
}
