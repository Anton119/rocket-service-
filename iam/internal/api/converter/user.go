package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	"github.com/Anton119/rocket-service-/iam/internal/service/input"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	commonv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/common/v1"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

// ParseUUID разбирает строковый UUID из gRPC-запроса.
func ParseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errs.ErrEmptySessionID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return id, nil
}

// ParseUserUUID разбирает user_uuid из gRPC-запроса.
func ParseUserUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return id, nil
}

// RegisterInputFromRequest собирает вход use case из rpc Register.
func RegisterInputFromRequest(req *userv1.RegisterRequest) input.RegisterInput {
	info := req.GetInfo()
	if info == nil {
		return input.RegisterInput{}
	}

	userInfo := info.GetInfo()
	if userInfo == nil {
		return input.RegisterInput{
			Password: info.GetPassword(),
		}
	}

	return input.RegisterInput{
		Login:    userInfo.GetLogin(),
		Email:    userInfo.GetEmail(),
		Password: info.GetPassword(),
	}
}

// LoginInputFromRequest собирает вход use case из rpc Login.
func LoginInputFromRequest(req *authv1.LoginRequest) input.LoginInput {
	return input.LoginInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
}

// SessionToProto переводит доменную сессию в protobuf Session.
func SessionToProto(s model.Session) *commonv1.Session {
	msg := &commonv1.Session{
		Uuid:      s.UUID.String(),
		CreatedAt: timestamppb.New(s.CreatedAt),
		ExpiresAt: timestamppb.New(s.ExpiresAt),
	}

	if s.UpdatedAt != nil {
		msg.UpdatedAt = timestamppb.New(*s.UpdatedAt)
	}

	return msg
}

// UserToProto переводит доменного пользователя в protobuf User.
func UserToProto(u model.User, email string) *commonv1.User {
	msg := &commonv1.User{
		Uuid: u.UUID.String(),
		Info: &commonv1.UserInfo{
			Login: u.Login,
			Email: email,
		},
		CreatedAt: timestamppb.New(u.CreatedAt),
	}

	if u.UpdatedAt != nil {
		msg.UpdatedAt = timestamppb.New(*u.UpdatedAt)
	}

	return msg
}
