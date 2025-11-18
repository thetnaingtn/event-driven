package event

import (
	"context"
	"tickets/entity"
)

func (h Handler) CancelBooking(ctx context.Context, event entity.TicketBookingCanceled) error {
	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-refund", []string{
		event.TicketID,
		event.CustomerEmail,
		event.Price.Amount,
		event.Price.Currency,
	})
}
