package adapters

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entity"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/dead_nation"
)

type BookingAPIClient struct {
	clients *clients.Clients
}

func NewBookingAPIClient(apiClients *clients.Clients) *BookingAPIClient {
	if apiClients == nil {
		panic("NewBookingAPIClient: clients is nil")
	}

	return &BookingAPIClient{
		clients: apiClients,
	}
}

func (c *BookingAPIClient) MakeBooking(ctx context.Context, booking entity.CreateBookingRequest) error {
	resp, err := c.clients.DeadNation.PostTicketBookingWithResponse(ctx, dead_nation.PostTicketBookingRequest{
		EventId:         booking.EventID,
		CustomerAddress: booking.CustomerEmail,
		NumberOfTickets: booking.NumberOfTickets,
		BookingId:       booking.BookingID,
	})

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code while making booking %s: %d", booking.BookingID, resp.StatusCode())
	}

	return nil
}
