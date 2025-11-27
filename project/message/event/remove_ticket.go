package event

import (
	"context"
	"tickets/entity"
)

func (h *Handler) RemoveTicket(ctx context.Context, event *entity.TicketBookingCanceled) error {
	return h.ticketRepository.RemoveTicket(ctx, event.TicketID)
}
