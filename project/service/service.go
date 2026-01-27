package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	stdHTTP "net/http"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"tickets/db"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/command"
	"tickets/message/event"
	"tickets/message/outbox"
	"tickets/migrations"
)

type Service struct {
	db                 *sqlx.DB
	watermillRouter    *watermillMessage.Router
	echoRouter         *echo.Echo
	dataLakeRepository db.DataLake
	readModel          db.OpsBookingReadModel
}

type ReceiptService interface {
	event.ReceiptsService
	command.ReceiptsService
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	deadNationAPI event.DeadNationAPI,
	spreadsheetsAPI event.SpreadsheetsAPI,
	receiptsService ReceiptService,
	filesAPI event.FilesAPI,
	paymentsService command.PaymentsService,
) Service {
	watermillLogger := watermill.NewSlogLogger(log.FromContext(context.Background()))

	redisPublisher := message.NewRedisPublisher(redisClient, watermillLogger)
	redisSubscriber := message.NewRedisSubscriber(redisClient, watermillLogger)

	eventBus := event.NewBus(redisPublisher)

	ticketsRepo := db.NewTicketsRepository(dbConn)
	OpsBookingReadModel := db.NewOpsBookingReadModel(dbConn, eventBus)
	showsRepo := db.NewShowsRepository(dbConn)
	bookingsRepository := db.NewBookingsRepository(dbConn)
	dataLakeRepository := db.NewDataLake(dbConn)

	eventsHandler := event.NewHandler(
		deadNationAPI,
		spreadsheetsAPI,
		receiptsService,
		filesAPI,
		ticketsRepo,
		showsRepo,
		eventBus,
		dataLakeRepository,
	)

	commandsHandler := command.NewHandler(
		eventBus,
		receiptsService,
		paymentsService,
	)

	commandBus := command.NewBus(redisPublisher, command.NewBusConfig(watermillLogger))

	postgresSubscriber := outbox.NewPostgresSubscriber(dbConn.DB, watermillLogger)
	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)
	commandProcessorConfig := command.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		postgresSubscriber,
		redisPublisher,
		redisSubscriber,
		eventProcessorConfig,
		eventsHandler,
		commandProcessorConfig,
		commandsHandler,
		OpsBookingReadModel,
		watermillLogger,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
		commandBus,
		ticketsRepo,
		showsRepo,
		bookingsRepository,
		OpsBookingReadModel,
	)

	return Service{
		dbConn,
		watermillRouter,
		echoRouter,
		dataLakeRepository,
		OpsBookingReadModel,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDatabaseSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	errgrp.Go(func() error {
		for {
			events, err := s.dataLakeRepository.GetEvents(ctx)
			if err != nil && !errors.Is(sql.ErrNoRows, err) {
				return err
			}

			if len(events) > 0 {
				for _, event := range events {
					if err := migrations.MigrateEvent(ctx, event, s.readModel); err != nil {
						return err
					}
				}

				return nil
			}

			t := time.NewTimer(time.Millisecond * 500)

			select {
			case <-ctx.Done():
				t.Stop()
				return errors.New("canceled")
			case <-t.C:
				fmt.Println("retry event migration")
			}
		}
	})

	return errgrp.Wait()
}
