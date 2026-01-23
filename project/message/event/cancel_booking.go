package event

import (
	"context"
	"errors"
	"tickets/entities"
)

func (h Handler) CancelBooking(ctx context.Context, event *entities.TicketBookingCanceled_v1) error {
	if event == nil {
		return errors.New("empty event received")
	}

	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-refund", []string{
		event.TicketID,
		event.CustomerEmail,
		event.Price.Amount,
		event.Price.Currency,
	})
}
