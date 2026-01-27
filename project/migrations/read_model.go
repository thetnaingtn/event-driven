package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"tickets/db"
	"tickets/entities"
	"time"

	"github.com/google/uuid"
)

type bookingMade_v0 struct {
	Header entities.MessageHeader `json:"header"`

	NumberOfTickets int `json:"number_of_tickets"`

	BookingID uuid.UUID `json:"booking_id"`

	CustomerEmail string    `json:"customer_email"`
	ShowID        uuid.UUID `json:"show_id"`
}

type ticketBookingConfirmed_v0 struct {
	Header entities.MessageHeader `json:"header"`

	TicketID      string         `json:"ticket_id"`
	CustomerEmail string         `json:"customer_email"`
	Price         entities.Money `json:"price"`

	BookingID string `json:"booking_id"`
}

type ticketReceiptIssued_v0 struct {
	Header entities.MessageHeader `json:"header"`

	TicketID      string `json:"ticket_id"`
	ReceiptNumber string `json:"receipt_number"`

	IssuedAt time.Time `json:"issued_at"`
}

type ticketPrinted_v0 struct {
	Header entities.MessageHeader `json:"header"`

	TicketID string `json:"ticket_id"`
	FileName string `json:"file_name"`
}

type ticketRefunded_v0 struct {
	Header entities.MessageHeader `json:"header"`

	TicketID string `json:"ticket_id"`
}

func unmarshalDataLakeEvent[T any](event entities.DataLakeEvent) (*T, error) {
	eventInstance := new(T)

	if err := json.Unmarshal(event.EventPayload, eventInstance); err != nil {
		return nil, fmt.Errorf("could not unmarshal event %s: %w", event.EventName, err)
	}

	return eventInstance, nil
}

func MigrateEvent(ctx context.Context, event entities.DataLakeEvent, rm db.OpsBookingReadModel) error {
	switch event.EventName {
	case "BookingMade_v0":
		bookingMade, err := unmarshalDataLakeEvent[bookingMade_v0](event)
		if err != nil {
			return err
		}

		return rm.OnBookingMade(ctx, &entities.BookingMade_v1{
			// you should map v0 event to your v1 event here
			Header:          bookingMade.Header,
			NumberOfTickets: bookingMade.NumberOfTickets,
			BookingID:       bookingMade.BookingID,
			CustomerEmail:   bookingMade.CustomerEmail,
			ShowID:          bookingMade.ShowID,
		})
	case "TicketBookingConfirmed_v0":
		bookingConfirmedEvent, err := unmarshalDataLakeEvent[ticketBookingConfirmed_v0](event)
		if err != nil {
			return err
		}

		return rm.OnTicketBookingConfirmed(ctx, &entities.TicketBookingConfirmed_v1{
			// you should map v0 event to your v1 event here
			Header:        bookingConfirmedEvent.Header,
			TicketID:      bookingConfirmedEvent.TicketID,
			CustomerEmail: bookingConfirmedEvent.CustomerEmail,
			Price:         bookingConfirmedEvent.Price,
			BookingID:     bookingConfirmedEvent.BookingID,
		})
	case "TicketReceiptIssued_v0":
		receiptIssuedEvent, err := unmarshalDataLakeEvent[ticketReceiptIssued_v0](event)
		if err != nil {
			return err
		}

		return rm.OnTicketReceiptIssued(ctx, &entities.TicketReceiptIssued_v1{
			// you should map v0 event to your v1 event here
			Header:        receiptIssuedEvent.Header,
			TicketID:      receiptIssuedEvent.TicketID,
			ReceiptNumber: receiptIssuedEvent.ReceiptNumber,
			IssuedAt:      receiptIssuedEvent.IssuedAt,
		})
	case "TicketPrinted_v0":
		ticketPrintedEvent, err := unmarshalDataLakeEvent[ticketPrinted_v0](event)
		if err != nil {
			return err
		}

		return rm.OnTicketPrinted(ctx, &entities.TicketPrinted_v1{
			// you should map v0 event to your v1 event here
			Header:   ticketPrintedEvent.Header,
			TicketID: ticketPrintedEvent.TicketID,
			FileName: ticketPrintedEvent.FileName,
		})
	case "TicketRefunded_v0":
		ticketRefundedEvent, err := unmarshalDataLakeEvent[ticketRefunded_v0](event)
		if err != nil {
			return err
		}

		return rm.OnTicketRefunded(ctx, &entities.TicketRefunded_v1{
			// you should map v0 event to your v1 event here
			Header:   ticketRefundedEvent.Header,
			TicketID: ticketRefundedEvent.TicketID,
		})
	default:
		return fmt.Errorf("unknown event %s", event.EventName)
	}
}
