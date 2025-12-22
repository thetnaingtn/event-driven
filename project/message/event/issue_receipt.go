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
		TicketID:       event.TicketID,
		Price:          event.Price,
		IdempotencyKey: event.Header.IdempotencyKey,
	}

	resp, err := h.receiptService.IssueReceipt(ctx, request)
	if err != nil {
		return err
	}

	return h.eventBus.Publish(ctx, entity.TicketReceiptIssued{
		Header:        entity.NewMessageHeaderWithIdempotencyKey(*event.Header.IdempotencyKey),
		TicketID:      event.TicketID,
		ReceiptNumber: resp.ReceiptNumber,
		IssuedAt:      resp.IssuedAt,
	})
}
