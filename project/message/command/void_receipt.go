package command

import (
	"context"
	"tickets/entity"
)

func (h Handler) VoidReceipt(ctx context.Context, command *entity.RefundTicket) error {
	return h.receiptClient.VoidReceipt(ctx, entity.VoidReceiptRequest{
		TicketId:       command.TicketID,
		IdempotencyKey: command.Header.IdempotencyKey,
	})
}
