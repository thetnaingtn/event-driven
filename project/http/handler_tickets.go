package http

import (
	"log/slog"
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
		msg := []byte(ticket)
		err := h.publisher.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), msg))
		if err != nil {
			slog.With("error", err).Error("Error publishing issue receipt")
		}

		err = h.publisher.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), msg))
		if err != nil {
			slog.With("error", err).Error("Error publishing append to tracker")
		}
	}

	return c.NoContent(http.StatusOK)
}
