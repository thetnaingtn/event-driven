package command

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewBus(pub message.Publisher, config cqrs.CommandBusConfig) *cqrs.CommandBus {
	bus, err := cqrs.NewCommandBusWithConfig(pub, config)

	if err != nil {
		panic(fmt.Errorf("can't create command bus: %w", err))
	}

	return bus
}
