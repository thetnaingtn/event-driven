package message

import (
	"log/slog"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func useMiddleware(router *message.Router) {
	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(addCorrelationID)
	router.AddMiddleware(logger)
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
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context())
		logger.With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		).Info("Handling a message")
		return next(msg)
	}
}
