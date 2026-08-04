package kafka

import (
	"context"

	"github.com/Anton119/rocket-service-/platform/pkg/auth"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
)

// ProducerSessionHeaders собирает Kafka headers с session-uuid из context.
func ProducerSessionHeaders(ctx context.Context) []kafka.Header {
	sessionUUID, ok := auth.SessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return nil
	}

	return []kafka.Header{{
		Key:   auth.SessionMetadataKey,
		Value: []byte(sessionUUID),
	}}
}

// ConsumerSession — middleware, кладущий session-uuid из Kafka headers в context.
func ConsumerSession() kafka.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			for _, h := range msg.Headers {
				if h.Key == auth.SessionMetadataKey && len(h.Value) > 0 {
					ctx = auth.WithSessionUUID(ctx, string(h.Value))
					break
				}
			}

			return next(ctx, msg)
		}
	}
}
