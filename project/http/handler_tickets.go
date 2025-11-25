package http

import (
	"log/slog"
	"net/http"
	"tickets/entity"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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

	bus, err := cqrs.NewEventBusWithConfig(log.CorrelationPublisherDecorator{Publisher: h.publisher}, cqrs.EventBusConfig{
		GeneratePublishTopic: func(geptp cqrs.GenerateEventPublishTopicParams) (string, error) {
			return geptp.EventName, nil
		},
		Marshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
	})

	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		switch ticket.Status {
		case "confirmed":
			bookingConfirmedEvent := entity.TicketBookingConfirmed{
				Header:        entity.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			slog.Info("Publishing ticket booking confirmed event")

			if err := bus.Publish(c.Request().Context(), bookingConfirmedEvent); err != nil {
				return err
			}

		case "canceled":
			bookingCanceledEvent := entity.TicketBookingCanceled{
				Header:        entity.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			slog.Info("Publishing ticket booking canceled event")

			if err := bus.Publish(c.Request().Context(), bookingCanceledEvent); err != nil {
				return err
			}
		}

	}

	return c.NoContent(http.StatusOK)
}
