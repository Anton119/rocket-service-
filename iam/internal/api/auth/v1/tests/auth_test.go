package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authapi "github.com/Anton119/rocket-service-/iam/internal/api/auth/v1"
	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
)

type authServiceMock struct {
	mock.Mock
}

func (m *authServiceMock) Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	args := m.Called(ctx, in)
	id, _ := args.Get(0).(uuid.UUID)
	return id, args.Error(1)
}

func (m *authServiceMock) Whoami(ctx context.Context, sessionUUID uuid.UUID) (model.Session, model.User, error) {
	args := m.Called(ctx, sessionUUID)
	session, _ := args.Get(0).(model.Session)
	user, _ := args.Get(1).(model.User)
	return session, user, args.Error(2)
}

func (m *authServiceMock) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	args := m.Called(ctx, sessionUUID)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	sessionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")

	tests := []struct {
		name      string
		req       *authv1.LoginRequest
		setupMock func(m *authServiceMock)
		wantErr   error
	}{
		{
			name: "успешный вход",
			req: &authv1.LoginRequest{
				Login:    "alice",
				Password: "password123",
			},
			setupMock: func(m *authServiceMock) {
				m.On("Login", ctx, input.LoginInput{Login: "alice", Password: "password123"}).
					Return(sessionUUID, nil)
			},
		},
		{
			name: "неверные credentials",
			req: &authv1.LoginRequest{
				Login:    "alice",
				Password: "wrong",
			},
			setupMock: func(m *authServiceMock) {
				m.On("Login", ctx, mock.Anything).
					Return(uuid.Nil, errs.ErrInvalidCredentials)
			},
			wantErr: errs.ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &authServiceMock{}
			tc.setupMock(svc)

			resp, err := authapi.NewAPI(svc).Login(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, sessionUUID.String(), resp.GetSessionUuid())
		})
	}
}

func TestWhoami(t *testing.T) {
	ctx := context.Background()
	sessionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
	userUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
	now := time.Now()

	tests := []struct {
		name      string
		req       *authv1.WhoamiRequest
		setupMock func(m *authServiceMock)
		wantErr   error
	}{
		{
			name: "успешный whoami",
			req:  &authv1.WhoamiRequest{SessionUuid: sessionUUID.String()},
			setupMock: func(m *authServiceMock) {
				m.On("Whoami", ctx, sessionUUID).
					Return(model.Session{UUID: sessionUUID, CreatedAt: now, ExpiresAt: now.Add(time.Hour)}, model.User{
						UUID:      userUUID,
						Login:     "alice",
						CreatedAt: now,
					}, nil)
			},
		},
		{
			name:    "неверный session_uuid",
			req:     &authv1.WhoamiRequest{SessionUuid: "not-a-uuid"},
			wantErr: errs.ErrInvalidUUID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &authServiceMock{}
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			resp, err := authapi.NewAPI(svc).Whoami(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, sessionUUID.String(), resp.GetSession().GetUuid())
			assert.Equal(t, userUUID.String(), resp.GetUser().GetUuid())
			assert.Equal(t, "alice", resp.GetUser().GetInfo().GetLogin())
		})
	}
}

func TestLogout(t *testing.T) {
	ctx := context.Background()
	sessionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")

	svc := &authServiceMock{}
	svc.On("Logout", ctx, sessionUUID).Return(nil)

	resp, err := authapi.NewAPI(svc).Logout(ctx, &authv1.LogoutRequest{
		SessionUuid: sessionUUID.String(),
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
