package command

import (
	"context"
	"tickets/entity"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ReceiptClient interface {
	VoidReceipt(ctx context.Context, request entity.VoidReceipt) error
}

type PaymentClient interface {
	Refund(ctx context.Context, request entity.RefundTicketRequest) error
}

type Handler struct {
	receiptClient ReceiptClient
	paymentClient PaymentClient
	eventBus      *cqrs.EventBus
}

func NewHandler(receiptClient ReceiptClient, paymentClient PaymentClient, eventBus *cqrs.EventBus) Handler {
	return Handler{
		receiptClient: receiptClient,
		paymentClient: paymentClient,
		eventBus:      eventBus,
	}
}
