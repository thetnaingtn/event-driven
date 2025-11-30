package http

import (
	"net/http"
	"tickets/entity"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookingRequest struct {
	ShowID          string `json:"show_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}

func (h *Handler) BookTickets(c echo.Context) error {
	var request BookingRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	booking := &entity.Booking{
		BookingID:       uuid.NewString(),
		ShowID:          request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	}

	if err := h.bookingRepository.CreateBooking(c.Request().Context(), booking); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		BookingID string `json:"booking_id"`
	}{
		BookingID: booking.BookingID,
	})
}
