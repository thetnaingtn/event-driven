package event

import (
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewBus(pub message.Publisher) *cqrs.EventBus {
	eventBus, err := cqrs.NewEventBusWithConfig(
		pub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				event, ok := params.Event.(entities.Event)
				if !ok {
					return "", fmt.Errorf("invalid event type: %T doesn't implement entities.Event", params.Event)
				}

				if event.IsInternal() {
					return "internal-events.svc-tickets." + params.EventName, nil
				} else {
					return "events", nil
				}
			},
			Marshaler: marshaler,
		},
	)
	if err != nil {
		panic(err)
	}

	return eventBus
}
