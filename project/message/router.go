package message

import (
	"encoding/json"
	"log/slog"
	"tickets/entity"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
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

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)

	if err != nil {
		panic(err)
	}

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)

	if err != nil {
		panic(err)
	}

	appendToRefundSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-refund",
	}, logger)

	if err != nil {
		panic(err)
	}

	handler := event.NewHandler(spreadsheetClient, receiptClient)

	router.AddConsumerHandler("tracker-handler", "TicketBookingConfirmed", appendToTrackerSub, func(msg *message.Message) error {
		var event entity.TicketBookingConfirmed
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}

		slog.Info("Appending ticket to tracker")

		if err := handler.AddTracker(msg.Context(), event); err != nil {
			return err
		}

		return nil
	})

	router.AddConsumerHandler("issue-receipt-handler", "TicketBookingConfirmed", issueReceiptSub, func(msg *message.Message) error {
		var event entity.TicketBookingConfirmed
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		if event.Price.Currency == "" {
			event.Price.Currency = "USD"
		}

		slog.Info("Issuing receipt")

		if err := handler.IssueReceipt(msg.Context(), event); err != nil {
			return err
		}

		return nil
	})

	router.AddConsumerHandler("cancel_ticket", "TicketBookingCanceled", appendToRefundSub, func(msg *message.Message) error {
		if msg.Metadata.Get("type") != "TicketBookingCanceled" {
			slog.Error("Invalid message type")
			return nil
		}

		var event entity.TicketBookingCanceled
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		slog.Info("Appending cancellation")

		return handler.CancelBooking(msg.Context(), event)
	})

	return router
}
