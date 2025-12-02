package message

import (
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewRouter(
	postgresSubscriber *watermillSQL.Subscriber,
	publisher message.Publisher,
	config cqrs.EventProcessorConfig,
	logger watermill.LoggerAdapter,
	eventHandler event.Handler,
) *message.Router {
	router := message.NewDefaultRouter(logger)
	useMiddleware(router, logger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, config)

	outbox.AddForwardHandler(postgresSubscriber, publisher, logger, router)

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
		cqrs.NewEventHandler(
			"StoreTicketToDB",
			eventHandler.SaveTicket,
		),
		cqrs.NewEventHandler(
			"RemoveTicketFromDB",
			eventHandler.RemoveTicket,
		),
		cqrs.NewEventHandler(
			"PrintTicket",
			eventHandler.PrintTicket,
		),
	)

	return router
}
