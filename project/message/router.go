package message

import (
	"tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewRouter(
	config cqrs.EventProcessorConfig,
	logger watermill.LoggerAdapter,
	eventHandler event.Handler,
) *message.Router {
	router := message.NewDefaultRouter(logger)
	useMiddleware(router, logger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, config)

	if err != nil {
		panic(err)
	}

	eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"AppendToTracker",
			eventHandler.AppendToTracker,
		),
		cqrs.NewEventHandler(
			"TicketRefundToSheet",
			eventHandler.TicketRefundToSheet,
		),
		cqrs.NewEventHandler(
			"IssueReceipt",
			eventHandler.IssueReceipt,
		),
	)

	return router
}
