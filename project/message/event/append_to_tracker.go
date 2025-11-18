package event

import (
	"context"
	"tickets/entity"
)

func (h Handler) AddTracker(ctx context.Context, event entity.TicketBookingConfirmed) error {
	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{
		event.TicketID,
		event.CustomerEmail,
		event.Price.Amount,
		event.Price.Currency,
	})
}
