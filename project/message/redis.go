package message

import (
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewPublisher(rdb *redis.Client, logger watermill.LoggerAdapter) (message.Publisher, error) {
	var pub message.Publisher
	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		return nil, err
	}

	pub = log.CorrelationPublisherDecorator{Publisher: pub}

	return pub, err
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
