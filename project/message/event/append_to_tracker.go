package event

import (
	"context"
	"errors"
	"tickets/entity"
)

func (h Handler) AddTracker(ctx context.Context, event *entity.TicketBookingConfirmed) error {
	if event == nil {
		return errors.New("empty event received")
	}

	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{
		event.TicketID,
		event.CustomerEmail,
		event.Price.Amount,
		event.Price.Currency,
	})
}
