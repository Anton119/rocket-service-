package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
)

// Client — gRPC-клиент IAM AuthService.
type Client struct {
	api authv1.AuthServiceClient
}

// NewClient создаёт клиент аутентификации.
func NewClient(api authv1.AuthServiceClient) *Client {
	return &Client{api: api}
}

// Whoami проверяет сессию и возвращает UUID пользователя.
func (c *Client) Whoami(ctx context.Context, sessionUUID string) (string, error) {
	resp, err := c.api.Whoami(ctx, &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			return "", fmt.Errorf("сессия недействительна: %w", err)
		}

		return "", fmt.Errorf("whoami: %w", err)
	}

	return resp.GetUser().GetUuid(), nil
}
