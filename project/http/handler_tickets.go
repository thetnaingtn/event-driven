package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"tickets/entity"

	"github.com/google/uuid"
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

	idempotencyKey := c.Request().Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	for _, ticket := range request.Tickets {
		messageHeader := entity.NewMessageHeaderWithIdempotencyKey(idempotencyKey + ticket.TicketID)

		switch ticket.Status {
		case "confirmed":
			bookingConfirmedEvent := entity.TicketBookingConfirmed{
				Header:        messageHeader,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			slog.Info("Publishing ticket booking confirmed event")

			if err := h.eventBus.Publish(c.Request().Context(), bookingConfirmedEvent); err != nil {
				return err
			}

		case "canceled":
			bookingCanceledEvent := entity.TicketBookingCanceled{
				Header:        messageHeader,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			slog.Info("Publishing ticket booking canceled event")

			if err := h.eventBus.Publish(c.Request().Context(), bookingCanceledEvent); err != nil {
				return err
			}
		}

	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) GetAllTickets(c echo.Context) error {
	tickets, err := h.ticketRepository.FindAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tickets)
}

func (h Handler) PutTicketRefund(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	slog.Info("Sending ticket refund command")

	if err := h.commandBus.Send(c.Request().Context(), entity.RefundTicket{
		TicketID: ticketID,
		Header:   entity.NewMessageHeaderWithIdempotencyKey(uuid.NewString()),
	}); err != nil {
		return fmt.Errorf("can't publish")
	}

	return c.NoContent(http.StatusAccepted)
}
