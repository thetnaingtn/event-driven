package http

import (
	"encoding/json"
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
		receiptPayload := entity.IssueReceiptPayload{
			TicketID: ticket.TicketID,
			Price:    ticket.Price,
		}

		bytes, err := json.Marshal(receiptPayload)
		if err != nil {
			slog.With("error", err).Error("Error marshaling receipt payload")
			return err
		}

		err = h.publisher.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), bytes))
		if err != nil {
			slog.With("error", err).Error("Error publishing issue receipt")
		}

		trackerPayload := entity.AppendToTrackerPayload{
			TicketID:      ticket.TicketID,
			CustomerEmail: ticket.CustomerEmail,
			Price:         ticket.Price,
		}

		bytes, err = json.Marshal(trackerPayload)
		if err != nil {
			slog.With("error", err).Error("Error marshaling tracker payload")
			return err
		}

		err = h.publisher.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), bytes))
		if err != nil {
			slog.With("error", err).Error("Error publishing append to tracker")
		}
	}

	return c.NoContent(http.StatusOK)
}
