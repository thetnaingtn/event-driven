package http

import (
	"net/http"
	"tickets/entity"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ShowRequest struct {
	DeadNationID    string    `json:"dead_nation_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	StartTime       time.Time `json:"start_time"`
	Title           string    `json:"title"`
	Venue           string    `json:"venue"`
}

func (h Handler) CreateShow(c echo.Context) error {
	var request ShowRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	showID := uuid.NewString()

	show := &entity.Show{
		ShowID:          showID,
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}

	if err := h.showRepository.CreateShow(c.Request().Context(), show); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		ShowID string `json:"show_id"`
	}{
		ShowID: showID,
	})
}
