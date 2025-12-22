package command

import (
	"context"
	"fmt"
	"tickets/entity"
)

func (h Handler) RefundTicket(ctx context.Context, command *entity.RefundTicket) error {
	var err error
	err = h.receiptClient.VoidReceipt(ctx, entity.RefundTicketRequest{
		TicketId:       command.TicketID,
		IdempotencyKey: command.Header.IdempotencyKey,
	})
	if err != nil {
		return err
	}

	err = h.paymentClient.Refund(ctx, entity.RefundTicketRequest{
		TicketId:       command.TicketID,
		IdempotencyKey: command.Header.IdempotencyKey,
		RefundReason:   "customer requested refund",
	})

	if err != nil {
		return err
	}

	err = h.eventBus.Publish(ctx, entity.TicketRefunded{
		TicketID: command.TicketID,
	})

	if err != nil {
		return fmt.Errorf("failed to publish TicketRefunded event: %w", err)
	}

	return nil
}
