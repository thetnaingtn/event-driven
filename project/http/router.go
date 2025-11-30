package http

import (
	"log"
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	ticketRepository ticketRepository,
	showRepository showRepository,
	bookingRepository bookingRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus:          eventBus,
		ticketRepository:  ticketRepository,
		showRepository:    showRepository,
		bookingRepository: bookingRepository,
	}

	e.POST("/tickets-status", handler.PostTicketsConfirmation)
	e.GET("/health", func(c echo.Context) error {
		log.Println("here")
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/tickets", handler.GetAllTickets)

	e.POST("/shows", handler.CreateShow)

	e.POST("/book-tickets", handler.BookTickets)

	return e
}
