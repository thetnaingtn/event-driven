package message

import (
	"context"
	"tickets/entity"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
	spreadsheetClient event.SpreadSheetClient,
	receiptClient event.ReceiptClient,
) *message.Router {
	router := message.NewDefaultRouter(logger)
	useMiddleware(router, logger)

	processor, err := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(epgstp cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return epgstp.EventName, nil
		},
		Marshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
		SubscriberConstructor: func(epscp cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			switch epscp.HandlerName {
			case "tracker-handler":
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "tracker-handler",
				}, logger)
			case "issue-receipt-handler":
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "issue-receipt-handler",
				}, logger)
			case "refund-handler":
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "refund-handler",
				}, logger)
			default:
				return nil, nil
			}
		},
	})

	if err != nil {
		panic(err)
	}

	handler := event.NewHandler(spreadsheetClient, receiptClient)

	trackerHandler := cqrs.NewEventHandler("tracker-handler", func(ctx context.Context, event *entity.TicketBookingConfirmed) error {
		return handler.AddTracker(ctx, event)
	})

	issueReceiptHandler := cqrs.NewEventHandler("issue-receipt-handler", func(ctx context.Context, event *entity.TicketBookingConfirmed) error {
		return handler.IssueReceipt(ctx, event)
	})

	appendToRefundHandler := cqrs.NewEventHandler("refund-handler", func(ctx context.Context, event *entity.TicketBookingCanceled) error {
		return handler.CancelBooking(ctx, event)
	})

	if err := processor.AddHandlers(trackerHandler, issueReceiptHandler, appendToRefundHandler); err != nil {
		panic(err)
	}

	return router
}
