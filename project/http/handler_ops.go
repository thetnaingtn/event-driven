package http

import (
	"fmt"
	"net/http"
	"tickets/entities"
	"time"

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
	receiptIssueDate := c.QueryParam("receipt_issue_date")

	var (
		bookings []entities.OpsBooking
		err      error
	)

	filter := entities.Filter{}

	if receiptIssueDate != "" {
		date, err := time.Parse("2006-01-02", receiptIssueDate)
		if err != nil {
			return fmt.Errorf("failed to parse receipt_issue_date: %w", err)
		}
		filter.ReceiptIssueDate = &date
	}

	bookings, err = h.opsBookingReadModel.AllBookings(c.Request().Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, bookings)
}
