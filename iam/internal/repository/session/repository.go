package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Anton119/rocket-service-/iam/internal/model"
	repoconv "github.com/Anton119/rocket-service-/iam/internal/repository/converter"
	"github.com/Anton119/rocket-service-/iam/internal/repository/redis_view"
)

const sessionKeyPrefix = "session:"

// Repository — Redis-хранилище сессий.
type Repository struct {
	client *redis.Client
}

// New создаёт репозиторий сессий.
func New(client *redis.Client) *Repository {
	return &Repository{client: client}
}

// Create сохраняет сессию в Redis как HashMap с TTL.
func (r *Repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	key := sessionKey(session.UUID)
	view := repoconv.SessionToRedisView(session)

	if err := r.client.HSet(ctx, key, view).Err(); err != nil {
		return fmt.Errorf("сохранить сессию: %w", err)
	}

	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("установить TTL сессии: %w", err)
	}

	return nil
}

// Get возвращает сессию по UUID или ErrSessionNotFound.
func (r *Repository) Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error) {
	key := sessionKey(sessionUUID)

	var view redis_view.SessionRedisView
	if err := r.client.HGetAll(ctx, key).Scan(&view); err != nil {
		return model.Session{}, fmt.Errorf("получить сессию: %w", err)
	}

	return repoconv.SessionFromRedisView(view)
}

// Delete удаляет сессию из Redis (идемпотентная операция).
func (r *Repository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	key := sessionKey(sessionUUID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("удалить сессию: %w", err)
	}

	return nil
}

func sessionKey(sessionUUID uuid.UUID) string {
	return sessionKeyPrefix + sessionUUID.String()
}
