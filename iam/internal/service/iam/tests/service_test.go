package iam_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	iamsvc "github.com/Anton119/rocket-service-/iam/internal/service/iam"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
)

const testSessionTTL = time.Hour

func newService() (*iamsvc.Service, *userRepoFake, *sessionRepoFake) {
	users := newUserRepoFake()
	sessions := newSessionRepoFake()
	return iamsvc.NewService(users, sessions, testSessionTTL), users, sessions
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := newService()

	_, err := svc.Register(ctx, input.RegisterInput{Login: "carol", Password: "password123"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		in      input.RegisterInput
		wantErr error
	}{
		{
			name: "успешная регистрация",
			in: input.RegisterInput{
				Login:    "alice",
				Password: "password123",
			},
		},
		{
			name: "пустой логин",
			in: input.RegisterInput{
				Login:    "  ",
				Password: "password123",
			},
			wantErr: errs.ErrInvalidLogin,
		},
		{
			name: "слабый пароль",
			in: input.RegisterInput{
				Login:    "bob",
				Password: "short",
			},
			wantErr: errs.ErrWeakPassword,
		},
		{
			name: "логин уже занят",
			in: input.RegisterInput{
				Login:    "carol",
				Password: "password123",
			},
			wantErr: errs.ErrUserAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := svc.Register(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, uuid.Nil, id)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, id)
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	svc, users, _ := newService()

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	require.NoError(t, err)

	userUUID := uuid.New()
	require.NoError(t, users.Create(ctx, model.User{
		UUID:         userUUID,
		Login:        "dave",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}))

	tests := []struct {
		name    string
		in      input.LoginInput
		wantErr error
	}{
		{
			name: "успешный вход",
			in: input.LoginInput{
				Login:    "dave",
				Password: "password123",
			},
		},
		{
			name:    "пустые credentials",
			in:      input.LoginInput{},
			wantErr: errs.ErrEmptyCredential,
		},
		{
			name: "неверный пароль",
			in: input.LoginInput{
				Login:    "dave",
				Password: "wrong",
			},
			wantErr: errs.ErrInvalidCredentials,
		},
		{
			name: "пользователь не найден",
			in: input.LoginInput{
				Login:    "ghost",
				Password: "password123",
			},
			wantErr: errs.ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sessionUUID, err := svc.Login(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, uuid.Nil, sessionUUID)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, sessionUUID)
		})
	}
}

func TestWhoami(t *testing.T) {
	ctx := context.Background()
	svc, users, sessions := newService()

	userUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
	sessionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
	now := time.Now()

	require.NoError(t, users.Create(ctx, model.User{
		UUID:         userUUID,
		Login:        "eve",
		PasswordHash: "hash",
		CreatedAt:    now,
	}))
	require.NoError(t, sessions.Create(ctx, model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "eve",
		CreatedAt: now,
		ExpiresAt: now.Add(testSessionTTL),
	}, testSessionTTL))

	t.Run("успешный whoami", func(t *testing.T) {
		session, user, err := svc.Whoami(ctx, sessionUUID)
		require.NoError(t, err)
		assert.Equal(t, sessionUUID, session.UUID)
		assert.Equal(t, userUUID, user.UUID)
		assert.Equal(t, "eve", user.Login)
		assert.Empty(t, user.PasswordHash)
	})

	t.Run("пустой session_uuid", func(t *testing.T) {
		_, _, err := svc.Whoami(ctx, uuid.Nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrEmptySessionID)
	})

	t.Run("сессия не найдена", func(t *testing.T) {
		_, _, err := svc.Whoami(ctx, uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrSessionNotFound)
	})
}

func TestLogout(t *testing.T) {
	ctx := context.Background()
	svc, _, sessions := newService()
	sessionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")

	require.NoError(t, sessions.Create(ctx, model.Session{
		UUID:      sessionUUID,
		UserUUID:  uuid.New(),
		Login:     "frank",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(testSessionTTL),
	}, testSessionTTL))

	t.Run("успешный logout", func(t *testing.T) {
		require.NoError(t, svc.Logout(ctx, sessionUUID))

		_, _, err := svc.Whoami(ctx, sessionUUID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrSessionNotFound)
	})

	t.Run("идемпотентный logout", func(t *testing.T) {
		require.NoError(t, svc.Logout(ctx, sessionUUID))
	})

	t.Run("пустой session_uuid", func(t *testing.T) {
		err := svc.Logout(ctx, uuid.Nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrEmptySessionID)
	})
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	svc, users, _ := newService()
	userUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")

	require.NoError(t, users.Create(ctx, model.User{
		UUID:         userUUID,
		Login:        "grace",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}))

	t.Run("успешное получение", func(t *testing.T) {
		user, err := svc.GetUser(ctx, userUUID)
		require.NoError(t, err)
		assert.Equal(t, userUUID, user.UUID)
		assert.Equal(t, "grace", user.Login)
		assert.Empty(t, user.PasswordHash)
	})

	t.Run("неверный uuid", func(t *testing.T) {
		_, err := svc.GetUser(ctx, uuid.Nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrInvalidUUID)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		_, err := svc.GetUser(ctx, uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUserNotFound)
	})
}
