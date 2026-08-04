package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	userapi "github.com/Anton119/rocket-service-/iam/internal/api/user/v1"
	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
	commonv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/common/v1"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

type userServiceMock struct {
	mock.Mock
}

func (m *userServiceMock) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	args := m.Called(ctx, in)
	id, _ := args.Get(0).(uuid.UUID)
	return id, args.Error(1)
}

func (m *userServiceMock) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	args := m.Called(ctx, userUUID)
	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func registerRequest(login, password string) *userv1.RegisterRequest {
	return &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: login,
			},
			Password: password,
		},
	}
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	userUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")

	tests := []struct {
		name      string
		req       *userv1.RegisterRequest
		setupMock func(m *userServiceMock)
		wantErr   error
	}{
		{
			name: "успешная регистрация",
			req:  registerRequest("alice", "password123"),
			setupMock: func(m *userServiceMock) {
				m.On("Register", ctx, input.RegisterInput{Login: "alice", Password: "password123"}).
					Return(userUUID, nil)
			},
		},
		{
			name: "слабый пароль",
			req:  registerRequest("bob", "short"),
			setupMock: func(m *userServiceMock) {
				m.On("Register", ctx, mock.Anything).
					Return(uuid.Nil, errs.ErrWeakPassword)
			},
			wantErr: errs.ErrWeakPassword,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &userServiceMock{}
			tc.setupMock(svc)

			resp, err := userapi.NewAPI(svc).Register(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, userUUID.String(), resp.GetUserUuid())
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	userUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
	now := time.Now()

	tests := []struct {
		name      string
		req       *userv1.GetUserRequest
		setupMock func(m *userServiceMock)
		wantErr   error
	}{
		{
			name: "успешное получение",
			req:  &userv1.GetUserRequest{UserUuid: userUUID.String()},
			setupMock: func(m *userServiceMock) {
				m.On("GetUser", ctx, userUUID).
					Return(model.User{
						UUID:      userUUID,
						Login:     "alice",
						CreatedAt: now,
					}, nil)
			},
		},
		{
			name:    "неверный user_uuid",
			req:     &userv1.GetUserRequest{UserUuid: "bad"},
			wantErr: errs.ErrInvalidUUID,
		},
		{
			name: "пользователь не найден",
			req:  &userv1.GetUserRequest{UserUuid: userUUID.String()},
			setupMock: func(m *userServiceMock) {
				m.On("GetUser", ctx, userUUID).
					Return(model.User{}, errs.ErrUserNotFound)
			},
			wantErr: errs.ErrUserNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &userServiceMock{}
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			resp, err := userapi.NewAPI(svc).GetUser(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, userUUID.String(), resp.GetUser().GetUuid())
			assert.Equal(t, "alice", resp.GetUser().GetInfo().GetLogin())
		})
	}
}
