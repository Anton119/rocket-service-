package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/iam/internal/api/converter"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

// Register реализует rpc Register.
func (a *API) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	userUUID, err := a.svc.Register(ctx, apiconv.RegisterInputFromRequest(req))
	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{
		UserUuid: userUUID.String(),
	}, nil
}

// GetUser реализует rpc GetUser.
func (a *API) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	userUUID, err := apiconv.ParseUserUUID(req.GetUserUuid())
	if err != nil {
		return nil, err
	}

	user, err := a.svc.GetUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: apiconv.UserToProto(user, ""),
	}, nil
}
