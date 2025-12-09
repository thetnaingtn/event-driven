package command

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

var (
	marshaler = cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
)

func getTopicName(commandName string) string {
	return fmt.Sprintf("commands.%s", commandName)
}

func NewBusConfig(logger watermill.LoggerAdapter) cqrs.CommandBusConfig {
	return cqrs.CommandBusConfig{
		GeneratePublishTopic: func(cbgptp cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return getTopicName(cbgptp.CommandName), nil
		},
		Marshaler: marshaler,
		Logger:    logger,
	}
}

func NewProcessorConfig(rdb *redis.Client, logger watermill.LoggerAdapter) cqrs.CommandProcessorConfig {
	return cqrs.CommandProcessorConfig{
		SubscriberConstructor: func(cpscp cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb,
				ConsumerGroup: "commands." + cpscp.HandlerName,
			}, logger)
		},
		GenerateSubscribeTopic: func(cpgstp cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
			return getTopicName(cpgstp.CommandName), nil
		},
		Marshaler: marshaler,
		Logger:    logger,
	}
}
