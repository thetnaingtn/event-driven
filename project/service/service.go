package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill"
	wMessage "github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/sync/errgroup"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"
)

type Service struct {
	echoRouter *echo.Echo
	router     *wMessage.Router
}

func New(
	rdb *redis.Client,
	spreadsheetsAPI event.SpreadSheetClient,
	receiptsService event.ReceiptClient,
) Service {
	logger := watermill.NewSlogLogger(slog.Default())

	eventHandler := event.NewHandler(spreadsheetsAPI, receiptsService)

	eventProcessorConfig := event.NewEventProcessorConfig(rdb, logger)

	router := message.NewRouter(eventProcessorConfig, logger, eventHandler)

	publisher, err := message.NewPublisher(rdb, logger)
	if err != nil {
		panic(err)
	}

	eventBus := event.NewBus(publisher)
	echoRouter := ticketsHttp.NewHttpRouter(eventBus)

	return Service{
		echoRouter: echoRouter,
		router:     router,
	}
}

func (s Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.router.Run(ctx); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-s.router.Running()
		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(ctx)
	})

	return g.Wait()
}
