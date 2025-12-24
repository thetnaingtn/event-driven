package adapters

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entity"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/receipts"
)

type ReceiptsServiceClient struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewReceiptsServiceClient(clients *clients.Clients) *ReceiptsServiceClient {
	if clients == nil {
		panic("NewReceiptsServiceClient: clients is nil")
	}

	return &ReceiptsServiceClient{clients: clients}
}

func (c ReceiptsServiceClient) IssueReceipt(ctx context.Context, request entity.IssueReceiptRequest) (*entity.IssueReceiptResponse, error) {
	resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, receipts.CreateReceipt{
		TicketId: request.TicketID,
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
		IdempotencyKey: request.IdempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to post receipt: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &entity.IssueReceiptResponse{
			ReceiptNumber: resp.JSON200.Number,
			IssuedAt:      resp.JSON200.IssuedAt,
		}, nil
	case http.StatusCreated:
		return &entity.IssueReceiptResponse{
			ReceiptNumber: resp.JSON201.Number,
			IssuedAt:      resp.JSON201.IssuedAt,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected status code for POST receipts-api/receipts: %d", resp.StatusCode())
	}
}

func (c ReceiptsServiceClient) VoidReceipt(ctx context.Context, req entity.VoidReceipt) error {
	_, err := c.clients.Receipts.PutVoidReceiptWithResponse(ctx, receipts.VoidReceiptRequest{
		IdempotentId: &req.IdempotencyKey,
		Reason:       req.Reason,
		TicketId:     req.TicketID,
	})

	if err != nil {
		return err
	}

	return nil
}
