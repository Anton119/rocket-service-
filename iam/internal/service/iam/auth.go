package iam

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
)

const bcryptCost = 10

func (s *Service) register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	login := strings.TrimSpace(in.Login)
	if login == "" {
		return uuid.Nil, errs.ErrInvalidLogin
	}

	if len(in.Password) < 8 {
		return uuid.Nil, errs.ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("хешировать пароль: %w", err)
	}

	userUUID := uuid.New()
	now := time.Now()

	user := model.User{
		UUID:         userUUID,
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    now,
	}

	if err = s.users.Create(ctx, user); err != nil {
		return uuid.Nil, err
	}

	slog.InfoContext(ctx, "пользователь зарегистрирован", "login", login, "user_uuid", userUUID)

	return userUUID, nil
}

func (s *Service) login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	login := strings.TrimSpace(in.Login)
	if login == "" || in.Password == "" {
		return uuid.Nil, errs.ErrEmptyCredential
	}

	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errs.IsUserNotFound(err) {
			return uuid.Nil, errs.ErrInvalidCredentials
		}

		return uuid.Nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return uuid.Nil, errs.ErrInvalidCredentials
	}

	now := time.Now()
	sessionUUID := uuid.New()

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}

	if err = s.sessions.Create(ctx, session, s.sessionTTL); err != nil {
		return uuid.Nil, fmt.Errorf("создать сессию: %w", err)
	}

	slog.InfoContext(ctx, "пользователь вошёл в систему", "login", login, "session_uuid", sessionUUID)

	return sessionUUID, nil
}

func (s *Service) whoami(ctx context.Context, sessionUUID uuid.UUID) (model.Session, model.User, error) {
	if sessionUUID == uuid.Nil {
		return model.Session{}, model.User{}, errs.ErrEmptySessionID
	}

	session, err := s.sessions.Get(ctx, sessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	user, err := s.users.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	user.PasswordHash = ""

	return session, user, nil
}

func (s *Service) logout(ctx context.Context, sessionUUID uuid.UUID) error {
	if sessionUUID == uuid.Nil {
		return errs.ErrEmptySessionID
	}

	if err := s.sessions.Delete(ctx, sessionUUID); err != nil {
		return err
	}

	slog.InfoContext(ctx, "сессия завершена", "session_uuid", sessionUUID)

	return nil
}

func (s *Service) getUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	if userUUID == uuid.Nil {
		return model.User{}, errs.ErrInvalidUUID
	}

	user, err := s.users.GetByUUID(ctx, userUUID)
	if err != nil {
		return model.User{}, err
	}

	user.PasswordHash = ""

	return user, nil
}
