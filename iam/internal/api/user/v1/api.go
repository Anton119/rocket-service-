package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

// UserService описывает сервисный слой пользователей для gRPC API.
type UserService interface {
	Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

// API — gRPC-адаптер UserService.
type API struct {
	userv1.UnimplementedUserServiceServer
	svc UserService
}

// NewAPI создаёт обработчик gRPC UserService.
func NewAPI(svc UserService) *API {
	return &API{svc: svc}
}
