package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	logger := watermill.NewSlogLogger(nil)

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: redisClient,
	}, logger)

	if err != nil {
		panic(err)
	}

	err = service.New(
		spreadsheetsAPI,
		receiptsService,
		publisher,
	).Run(context.Background())
	if err != nil {
		panic(err)
	}
}
