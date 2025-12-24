package command

import (
	"context"
	"fmt"
	"tickets/entity"
)

func (h Handler) RefundTicket(ctx context.Context, command *entity.RefundTicket) error {
	var err error
	err = h.receiptClient.VoidReceipt(ctx, entity.VoidReceipt{
		TicketID:       command.TicketID,
		IdempotencyKey: *command.Header.IdempotencyKey,
		Reason:         "ticket refunded",
	})
	if err != nil {
		return err
	}

	err = h.paymentClient.Refund(ctx, entity.RefundTicketRequest{
		TicketId:       command.TicketID,
		IdempotencyKey: command.Header.IdempotencyKey,
		RefundReason:   "ticket refunded",
	})

	if err != nil {
		return err
	}

	err = h.eventBus.Publish(ctx, entity.TicketRefunded{
		Header:   entity.NewMessageHeader(),
		TicketID: command.TicketID,
	})

	if err != nil {
		return fmt.Errorf("failed to publish TicketRefunded event: %w", err)
	}

	return nil
}
