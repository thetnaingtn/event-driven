package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	rdb *redis.Client,
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
) Service {
	logger := watermill.NewSlogLogger(nil)

	message.NewHandler(rdb, logger, spreadsheetsAPI, receiptsService)
	publisher, err := message.NewPublisher(rdb, logger)
	if err != nil {
		panic(err)
	}

	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
