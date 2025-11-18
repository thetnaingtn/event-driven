package event

import (
	"context"
	"tickets/entity"
)

type SpreadSheetClient interface {
	AppendRow(ctx context.Context, sheetName string, tickets []string) error
}

type ReceiptClient interface {
	IssueReceipt(ctx context.Context, request entity.IssueReceiptRequest) error
}

type Handler struct {
	spreadsheetsAPI SpreadSheetClient
	receiptService  ReceiptClient
}

func NewHandler(spreadsheetsAPI SpreadSheetClient, receiptsService ReceiptClient) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptService:  receiptsService,
	}
}
