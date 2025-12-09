package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill"
	wMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"tickets/db"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/command"
	"tickets/message/event"
	"tickets/message/outbox"
)

type Service struct {
	db         *sqlx.DB
	echoRouter *echo.Echo
	router     *wMessage.Router
}

type ReceiptService interface {
	event.ReceiptClient
	command.ReceiptClient
}

func New(
	dbConn *sqlx.DB,
	rdb *redis.Client,
	spreadsheetsAPI event.SpreadSheetClient,
	fileAPIClient event.FileAPIClient,
	bookingAPIClient event.BookingAPIClient,
	receiptService ReceiptService,
	paymentService command.PaymentClient,
) Service {
	logger := watermill.NewSlogLogger(slog.Default())

	publisher, err := message.NewPublisher(rdb, logger)
	if err != nil {
		panic(err)
	}

	eventBus := event.NewBus(publisher)
	commandBus := command.NewBus(publisher, command.NewBusConfig(logger))

	ticketRepository := db.NewTicketRepository(dbConn)
	showRepository := db.NewShowRepository(dbConn)
	bookingRepository := db.NewBookingRepository(dbConn)

	eventHandler := event.NewHandler(spreadsheetsAPI, receiptService, ticketRepository, fileAPIClient, eventBus, bookingAPIClient, showRepository)
	commandHandler := command.NewHandler(receiptService, paymentService)

	eventProcessorConfig := event.NewEventProcessorConfig(rdb, logger)
	commandProcessorConfig := command.NewProcessorConfig(rdb, logger)
	postgresSubscriber := outbox.NewPostgresSubscriber(dbConn, logger)
	router := message.NewRouter(postgresSubscriber, publisher, eventProcessorConfig, logger, eventHandler, commandProcessorConfig, commandHandler)

	echoRouter := ticketsHttp.NewHttpRouter(eventBus, ticketRepository, showRepository, bookingRepository, commandBus)

	return Service{
		db:         dbConn,
		echoRouter: echoRouter,
		router:     router,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDatabaseSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

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
