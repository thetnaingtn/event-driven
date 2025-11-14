package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
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
	spreadsheetClient SpreadSheetClient,
	receiptClient ReceiptClient,
) error {
	spreadsheetSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "spreadsheet",
	}, logger)
	if err != nil {
		return err
	}

	receiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "receipt",
	}, logger)
	if err != nil {
		return err
	}

	go func() {
		msgs, err := spreadsheetSub.Subscribe(context.Background(), "append-to-tracker")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := spreadsheetClient.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)}); err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	go func() {
		msgs, err := receiptSub.Subscribe(context.Background(), "issue-receipt")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := receiptClient.IssueReceipt(msg.Context(), string(msg.Payload)); err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	return nil
}
