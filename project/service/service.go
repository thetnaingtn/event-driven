package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill"
	wMessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter *echo.Echo
	router     *wMessage.Router
}

func New(
	rdb *redis.Client,
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
) Service {
	logger := watermill.NewSlogLogger(nil)
	router := wMessage.NewDefaultRouter(logger)

	message.NewHandler(rdb, logger, router, spreadsheetsAPI, receiptsService)
	publisher, err := message.NewPublisher(rdb, logger)
	if err != nil {
		panic(err)
	}

	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	return Service{
		echoRouter: echoRouter,
		router:     router,
	}
}

func (s Service) Run(ctx context.Context) error {
	go func() {
		if err := s.router.Run(context.Background()); err != nil {
			panic(err)
		}
	}()
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
