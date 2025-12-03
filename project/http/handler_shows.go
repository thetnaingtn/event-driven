package http

import (
	"net/http"
	"tickets/entity"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) CreateShow(c echo.Context) error {
	show := &entity.Show{}
	if err := c.Bind(show); err != nil {
		return err
	}

	show.ShowID = uuid.New()

	if err := h.showRepository.CreateShow(c.Request().Context(), show); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		ShowID string `json:"show_id"`
	}{
		ShowID: show.ShowID.String(),
	})
}
