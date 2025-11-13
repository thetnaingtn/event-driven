package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
	"tickets/worker"
)

type Service struct {
	echoRouter *echo.Echo
	worker     *worker.Worker
}

func New(
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
) Service {
	worker := worker.NewWorker(spreadsheetsAPI, receiptsService)
	echoRouter := ticketsHttp.NewHttpRouter(worker)

	return Service{
		echoRouter: echoRouter,
		worker:     worker,
	}
}

func (s Service) Run(ctx context.Context) error {
	go s.worker.Run(ctx)
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
