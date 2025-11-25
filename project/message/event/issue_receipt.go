package event

import (
	"context"
	"errors"
	"tickets/entity"
)

func (h Handler) IssueReceipt(ctx context.Context, event *entity.TicketBookingConfirmed) error {
	if event == nil {
		return errors.New("empty event received")
	}

	request := entity.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	}
	return h.receiptService.IssueReceipt(ctx, request)
}
