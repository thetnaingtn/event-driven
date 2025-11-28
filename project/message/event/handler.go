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

type TicketRepository interface {
	SaveTicket(ctx context.Context, ticket *entity.Ticket) (*entity.Ticket, error)
	RemoveTicket(ctx context.Context, id string) error
}

type FileAPIClient interface {
	UploadFile(ctx context.Context, fileName, fileContent string) error
}

type Handler struct {
	spreadsheetsAPI  SpreadSheetClient
	receiptService   ReceiptClient
	ticketRepository TicketRepository
	fileAPIClient    FileAPIClient
}

func NewHandler(spreadsheetsAPI SpreadSheetClient, receiptsService ReceiptClient, ticketRepository TicketRepository, fileAPIClient FileAPIClient) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		spreadsheetsAPI:  spreadsheetsAPI,
		receiptService:   receiptsService,
		ticketRepository: ticketRepository,
		fileAPIClient:    fileAPIClient,
	}
}
