package adapters

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entity"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/payments"
)

type PaymentServiceCliennt struct {
	clients *clients.Clients
}

func NewPaymentServiceClient(clients *clients.Clients) PaymentServiceCliennt {
	if clients == nil {
		panic("NewPaymentServiceClient: clients is nil")
	}

	return PaymentServiceCliennt{clients: clients}
}

func (p PaymentServiceCliennt) Refund(ctx context.Context, request entity.RefundTicketRequest) error {
	resp, err := p.clients.Payments.PutRefundsWithResponse(ctx, payments.PaymentRefundRequest{
		DeduplicationId:  request.IdempotencyKey,
		Reason:           request.RefundReason,
		PaymentReference: request.TicketId,
	})

	if err != nil {
		return fmt.Errorf("failed to post refund for payment %s: %w", request.TicketId, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected for /payments-api/refunds status code: %d", resp.StatusCode())
	}

	return nil
}
