package event

import (
	"context"
	"errors"
	"tickets/entity"
)

func (h Handler) CancelBooking(ctx context.Context, event *entity.TicketBookingCanceled) error {
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
