package command

import (
	"context"
	"tickets/entity"
)

type ReceiptClient interface {
	VoidReceipt(ctx context.Context, request entity.RefundTicketRequest) error
}

type PaymentClient interface {
	Refund(ctx context.Context, request entity.RefundTicketRequest) error
}

type Handler struct {
	receiptClient ReceiptClient
	paymentClient PaymentClient
}

func NewHandler(receiptClient ReceiptClient, paymentClient PaymentClient) Handler {
	return Handler{
		receiptClient: receiptClient,
		paymentClient: paymentClient,
	}
}
