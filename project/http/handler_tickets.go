package http

import (
	"log/slog"
	"net/http"
	"tickets/entity"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string       `json:"ticket_id"`
	Status        string       `json:"status"`
	Price         entity.Money `json:"price"`
	CustomerEmail string       `json:"customer_email"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		ticketId := ticket.TicketID
		err := h.publisher.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), message.Payload(ticketId)))
		if err != nil {
			slog.With("error", err).Error("Error publishing issue receipt")
		}

		err = h.publisher.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), message.Payload(ticketId)))
		if err != nil {
			slog.With("error", err).Error("Error publishing append to tracker")
		}
	}

	return c.NoContent(http.StatusOK)
}
