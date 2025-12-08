package command

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewBus(pub message.Publisher) *cqrs.CommandBus {
	bus, err := cqrs.NewCommandBusWithConfig(pub, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(cbgptp cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return cbgptp.CommandName, nil
		},
		Marshaler: marshaler,
	})

	if err != nil {
		panic(fmt.Errorf("can't create command bus: %w", err))
	}

	return bus
}
