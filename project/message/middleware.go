package message

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func useMiddleware(router *message.Router, watermillLogger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(addCorrelationID)
	router.AddMiddleware(logger)
	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermillLogger,
	}.Middleware)
}

func addCorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx := log.ContextWithCorrelationID(msg.Context(), msg.Metadata.Get("correlation_id"))
		ctx = log.ToContext(ctx, slog.With("correlation_id", log.CorrelationIDFromContext(ctx)))
		msg.SetContext(ctx)
		return next(msg)
	}
}

func logger(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		logger := log.FromContext(msg.Context())
		logger.With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		).Info("Handling a message")

		msgs, err = next(msg)

		if err != nil {
			logger.With("error", err, "message_id", msg.UUID).Error("Error while handling a message")
		}

		return
	}
}
