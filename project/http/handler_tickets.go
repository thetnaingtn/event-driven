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

	correlationID := c.Request().Header.Get("Correlation-ID")

	for _, ticket := range request.Tickets {
		switch ticket.Status {
		case "confirmed":
			bookingConfirmedEvent := entity.TicketBookingConfirmed{
				Header:        entity.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			jsonPayload, err := json.Marshal(bookingConfirmedEvent)
			if err != nil {
				slog.Error("Error when marshaling event")
				return err
			}

			slog.Info("Publishing ticket booking confirmed event")
			msg := message.NewMessage(watermill.NewUUID(), jsonPayload)

			msg.Metadata.Set("correlation_id", correlationID)
			msg.Metadata.Set("type", "TicketBookingConfirmed")

			h.publisher.Publish("TicketBookingConfirmed", msg)
		case "canceled":
			bookingCanceledEvent := entity.TicketBookingCanceled{
				Header:        entity.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			jsonPayload, err := json.Marshal(bookingCanceledEvent)
			if err != nil {
				slog.Error("Error when marshaling event")
				return err
			}

			slog.Info("Publishing ticket booking canceled event")
			msg := message.NewMessage(watermill.NewUUID(), jsonPayload)
			msg.Metadata.Set("correlation_id", correlationID)
			msg.Metadata.Set("type", "TicketBookingCanceled")

			h.publisher.Publish("TicketBookingCanceled", msg)
		}

	}

	return c.NoContent(http.StatusOK)
}
