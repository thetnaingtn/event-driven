package main

import (
	"context"
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		log.Fatal("can't start redis subscriber", err)
	}

	defer subscriber.Close()

	msgs, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		log.Fatal("can't subscribe to progress topic", err)
	}

	for msg := range msgs {
		log.Printf("Message ID: %s - %s", msg.UUID, msg.Payload)
		msg.Ack()
	}
}
