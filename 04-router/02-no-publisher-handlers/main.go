package main

import (
	"context"
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

	router := message.NewDefaultRouter(logger)

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		panic(err)
	}

	router.AddConsumerHandler("print-fahrenheit", "temperature-fahrenheit", sub, func(msg *message.Message) error {
		fahrenheit := string(msg.Payload)
		log.Printf("Temperature read: %s\n", fahrenheit)
		return nil
	})

	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
