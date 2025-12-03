package http

import (
	"net/http"
	"tickets/entity"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookingRequest struct {
	ShowID          uuid.UUID `json:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
}

func (h *Handler) BookTickets(c echo.Context) error {
	booking := &BookingRequest{}
	if err := c.Bind(booking); err != nil {
		return err
	}

	bookingID := uuid.New()

	if err := h.bookingRepository.CreateBooking(c.Request().Context(), &entity.Booking{
		BookingID:       bookingID,
		ShowID:          booking.ShowID,
		CustomerEmail:   booking.CustomerEmail,
		NumberOfTickets: booking.NumberOfTickets,
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		BookingID string `json:"booking_id"`
	}{
		BookingID: bookingID.String(),
	})
}
