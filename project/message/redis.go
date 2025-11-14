package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewPublisher(rdb *redis.Client, logger watermill.LoggerAdapter) (message.Publisher, error) {
	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		return nil, err
	}

	return pub, err
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
