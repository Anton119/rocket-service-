package testutil

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/metadata"

	"github.com/Anton119/rocket-service-/platform/pkg/auth"
)

const (
	// TestSessionUUID — фиксированная сессия для API-тестов.
	TestSessionUUID = "550e8400-e29b-41d4-a716-446655440099"
	// TestUserUUID — UUID пользователя, связанного с TestSessionUUID.
	TestUserUUID = "550e8400-e29b-41d4-a716-446655440010"
)

// AuthValidator — заглушка IAM Whoami для интеграционных тестов.
type AuthValidator struct{}

// Whoami возвращает TestUserUUID для TestSessionUUID.
func (AuthValidator) Whoami(_ context.Context, sessionUUID string) (string, error) {
	if sessionUUID == TestSessionUUID {
		return TestUserUUID, nil
	}

	return "", errors.New("сессия недействительна")
}

// SetAuthHeader добавляет Authorization: Bearer <session_uuid>.
func SetAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+TestSessionUUID)
}

// AuthGRPCContext добавляет session-uuid в outgoing gRPC metadata.
func AuthGRPCContext(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, auth.SessionMetadataKey, TestSessionUUID)
}
