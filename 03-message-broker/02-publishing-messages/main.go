package main

import (
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		log.Fatal("can't start redis publisher", err)
	}
	defer publisher.Close()

	publisher.Publish("progress",
		message.NewMessage(watermill.NewUUID(), []byte("50")),
		message.NewMessage(watermill.NewUUID(), []byte("100")),
	)

}
