package event

import (
	"context"
	"errors"
	"tickets/entity"
)

func (h *Handler) SaveTicket(ctx context.Context, event *entity.TicketBookingConfirmed) error {
	if event == nil {
		return errors.New("empty ticket")
	}

	ticket := &entity.Ticket{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	}

	if _, err := h.ticketRepository.SaveTicket(ctx, ticket); err != nil {
		return err
	}

	return nil
}
