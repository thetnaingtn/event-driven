package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type SpreadSheetClient interface {
	AppendRow(ctx context.Context, sheetName string, tickets []string) error
}

type ReceiptClient interface {
	IssueReceipt(ctx context.Context, ticket string) error
}

func NewHandler(
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
	router *message.Router,
	spreadsheetClient SpreadSheetClient,
	receiptClient ReceiptClient,
) error {
	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)
	if err != nil {
		return err
	}

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)
	if err != nil {
		return err
	}

	router.AddConsumerHandler("tracker-handler", "append-to-tracker", appendToTrackerSub, func(msg *message.Message) error {
		ticket := string(msg.Payload)

		if err := spreadsheetClient.AppendRow(msg.Context(), "tickets-to-print", []string{ticket}); err != nil {
			return err
		}
		return nil
	})

	router.AddConsumerHandler("issue-receipt-handler", "issue-receipt", issueReceiptSub, func(msg *message.Message) error {
		ticket := string(msg.Payload)
		if err := receiptClient.IssueReceipt(msg.Context(), ticket); err != nil {
			return err
		}

		return nil
	})

	return nil
}
