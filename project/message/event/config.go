package event

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

func NewEventProcessorConfig(rdb *redis.Client, logger watermill.LoggerAdapter) cqrs.EventProcessorConfig {
	return cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(epgstp cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return epgstp.EventName, nil
		},
		SubscriberConstructor: func(epscp cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb,
				ConsumerGroup: "svc-tickets." + epscp.HandlerName,
			}, logger)
		},
		Marshaler: marshaler,
		Logger:    logger,
	}
}
