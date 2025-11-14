package message

import (
	"context"
	"log/slog"

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

	go func() {
		msgs, err := appendToTrackerSub.Subscribe(context.Background(), "append-to-tracker")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := spreadsheetClient.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)}); err != nil {
				slog.With("error", err).Error("Error issuing receipt")
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	go func() {
		msgs, err := issueReceiptSub.Subscribe(context.Background(), "issue-receipt")
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			if err := receiptClient.IssueReceipt(msg.Context(), string(msg.Payload)); err != nil {
				slog.With("error", err).Error("Error appending to tracker")
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	return nil
}
