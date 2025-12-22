package event

import (
	"context"
	"tickets/entity"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type SpreadSheetClient interface {
	AppendRow(ctx context.Context, sheetName string, tickets []string) error
}

type ReceiptClient interface {
	IssueReceipt(ctx context.Context, request entity.IssueReceiptRequest) (*entity.IssueReceiptResponse, error)
}

type TicketRepository interface {
	SaveTicket(ctx context.Context, ticket *entity.Ticket) (*entity.Ticket, error)
	RemoveTicket(ctx context.Context, id string) error
}

type FileAPIClient interface {
	UploadFile(ctx context.Context, fileName, fileContent string) error
}

type BookingAPIClient interface {
	MakeBooking(ctx context.Context, request entity.CreateBookingRequest) error
}

type ShowRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Show, error)
}

type Handler struct {
	spreadsheetsAPI  SpreadSheetClient
	receiptService   ReceiptClient
	ticketRepository TicketRepository
	fileAPIClient    FileAPIClient
	eventBus         *cqrs.EventBus
	bookingAPIClient BookingAPIClient
	showRepository   ShowRepository
}

func NewHandler(spreadsheetsAPI SpreadSheetClient, receiptsService ReceiptClient, ticketRepository TicketRepository, fileAPIClient FileAPIClient, eventBus *cqrs.EventBus, bookingAPIClient BookingAPIClient, showRepository ShowRepository) Handler {
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
		eventBus:         eventBus,
		bookingAPIClient: bookingAPIClient,
		showRepository:   showRepository,
	}
}
