package http

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		h.publisher.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
		h.publisher.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
	}

	return c.NoContent(http.StatusOK)
}
