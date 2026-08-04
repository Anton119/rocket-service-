package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/repository/redis_view"
)

const redisTimeLayout = time.RFC3339Nano

// SessionToRedisView переводит доменную сессию в представление для Redis.
func SessionToRedisView(s model.Session) redis_view.SessionRedisView {
	return redis_view.SessionRedisView{
		UUID:      s.UUID.String(),
		UserUUID:  s.UserUUID.String(),
		Login:     s.Login,
		CreatedAt: s.CreatedAt.Format(redisTimeLayout),
		ExpiresAt: s.ExpiresAt.Format(redisTimeLayout),
	}
}

// SessionFromRedisView переводит Redis-представление в доменную модель.
func SessionFromRedisView(v redis_view.SessionRedisView) (model.Session, error) {
	if v.UUID == "" {
		return model.Session{}, errs.ErrSessionNotFound
	}

	sessionUUID, err := uuid.Parse(v.UUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("разобрать uuid сессии: %w", errs.ErrInvalidUUID)
	}

	userUUID, err := uuid.Parse(v.UserUUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("разобрать user_uuid: %w", errs.ErrInvalidUUID)
	}

	createdAt, err := time.Parse(redisTimeLayout, v.CreatedAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("разобрать created_at: %w", err)
	}

	expiresAt, err := time.Parse(redisTimeLayout, v.ExpiresAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("разобрать expires_at: %w", err)
	}

	return model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     v.Login,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}
