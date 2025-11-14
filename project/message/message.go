package message

import (
	"context"
	"os"

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

func Subscribe(ctx context.Context, spreadsheetClient SpreadSheetClient, receiptClient ReceiptClient) error {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

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
		msgs, err := spreadsheetSub.Subscribe(ctx, "append-to-tracker")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := spreadsheetClient.AppendRow(ctx, "tickets-to-print", []string{string(msg.Payload)}); err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	go func() {
		msgs, err := receiptSub.Subscribe(ctx, "issue-receipt")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := receiptClient.IssueReceipt(ctx, string(msg.Payload)); err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	return nil
}
