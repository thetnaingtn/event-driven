package command

import (
	"context"
	"tickets/entity"
)

type ReceiptClient interface {
	VoidReceipt(ctx context.Context, request entity.VoidReceiptRequest) error
}

type Handler struct {
	receiptClient ReceiptClient
}

func NewHandler(receiptClient ReceiptClient) Handler {
	return Handler{
		receiptClient: receiptClient,
	}
}
