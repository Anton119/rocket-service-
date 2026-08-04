package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
)

// AuthService описывает сервисный слой аутентификации для gRPC API.
type AuthService interface {
	Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionUUID uuid.UUID) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID uuid.UUID) error
}

// API — gRPC-адаптер AuthService.
type API struct {
	authv1.UnimplementedAuthServiceServer
	svc AuthService
}

// NewAPI создаёт обработчик gRPC AuthService.
func NewAPI(svc AuthService) *API {
	return &API{svc: svc}
}
