package iam

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
)

// UserRepository — доступ к пользователям в PostgreSQL.
type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

// SessionRepository — доступ к сессиям в Redis.
type SessionRepository interface {
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}

// Service — бизнес-логика IAM (регистрация, аутентификация, сессии).
type Service struct {
	users      UserRepository
	sessions   SessionRepository
	sessionTTL time.Duration
}

// NewService создаёт IAM-сервис.
func NewService(users UserRepository, sessions SessionRepository, sessionTTL time.Duration) *Service {
	return &Service{
		users:      users,
		sessions:   sessions,
		sessionTTL: sessionTTL,
	}
}

// Register регистрирует нового пользователя.
func (s *Service) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	return s.register(ctx, in)
}

// Login выполняет вход и создаёт сессию.
func (s *Service) Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	return s.login(ctx, in)
}

// Whoami возвращает сессию и пользователя по UUID сессии.
func (s *Service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (model.Session, model.User, error) {
	return s.whoami(ctx, sessionUUID)
}

// Logout завершает сессию (идемпотентная операция).
func (s *Service) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	return s.logout(ctx, sessionUUID)
}

// GetUser возвращает пользователя по UUID.
func (s *Service) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	return s.getUser(ctx, userUUID)
}
