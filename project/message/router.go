package message

import (
	"context"
	"encoding/json"
	"log/slog"
	"tickets/entity"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type SpreadSheetClient interface {
	AppendRow(ctx context.Context, sheetName string, tickets []string) error
}

type ReceiptClient interface {
	IssueReceipt(ctx context.Context, request entity.IssueReceiptRequest) error
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
		payload := &entity.AppendToTrackerPayload{}
		if err := json.Unmarshal(msg.Payload, payload); err != nil {
			return err
		}

		slog.Info("Appending ticket to tracker")

		if err := spreadsheetClient.AppendRow(msg.Context(), "tickets-to-print", []string{
			payload.TicketID,
			payload.CustomerEmail,
			payload.Price.Amount,
			payload.Price.Currency,
		}); err != nil {
			return err
		}
		return nil
	})

	router.AddConsumerHandler("issue-receipt-handler", "issue-receipt", issueReceiptSub, func(msg *message.Message) error {
		var payload entity.IssueReceiptPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		slog.Info("Issuing receipt")

		request := entity.IssueReceiptRequest(payload)

		if err := receiptClient.IssueReceipt(msg.Context(), request); err != nil {
			return err
		}

		return nil
	})

	return nil
}
