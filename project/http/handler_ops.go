package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetReadModel(c echo.Context) error {
	bookingID := c.Param("id")

	booking, err := h.opsBookingReadModel.GetReadModelByBookingID(c.Request().Context(), bookingID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, booking)
}

func (h Handler) AllBookings(c echo.Context) error {
	bookings, err := h.opsBookingReadModel.AllBookings(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, bookings)
}
