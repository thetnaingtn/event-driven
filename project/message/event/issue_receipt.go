package event

import (
	"context"
	"tickets/entity"
)

func (h Handler) IssueReceipt(ctx context.Context, event entity.TicketBookingConfirmed) error {
	request := entity.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	}
	return h.receiptService.IssueReceipt(ctx, request)
}
