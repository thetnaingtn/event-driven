package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter      *echo.Echo
	spreadsheetsApi ticketsHttp.SpreadsheetsAPI
	receiptsService ticketsHttp.ReceiptsService
}

func New(
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
	publisher *redisstream.Publisher,
) Service {
	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	return Service{
		echoRouter:      echoRouter,
		receiptsService: receiptsService,
		spreadsheetsApi: spreadsheetsAPI,
	}
}

func (s Service) Run(ctx context.Context) error {
	err := message.Subscribe(ctx, s.spreadsheetsApi, s.receiptsService)
	if err != nil {
		return err
	}

	err = s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
