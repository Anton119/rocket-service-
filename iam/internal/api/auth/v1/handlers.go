package v1

import (
	"context"

	apiconv "github.com/Anton119/rocket-service-/iam/internal/api/converter"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
)

// Login реализует rpc Login.
func (a *API) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	sessionUUID, err := a.svc.Login(ctx, apiconv.LoginInputFromRequest(req))
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID.String(),
	}, nil
}

// Whoami реализует rpc Whoami.
func (a *API) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	sessionUUID, err := apiconv.ParseUUID(req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	session, user, err := a.svc.Whoami(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	return &authv1.WhoamiResponse{
		Session: apiconv.SessionToProto(session),
		User:    apiconv.UserToProto(user, ""),
	}, nil
}

// Logout реализует rpc Logout.
func (a *API) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	sessionUUID, err := apiconv.ParseUUID(req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	if err = a.svc.Logout(ctx, sessionUUID); err != nil {
		return nil, err
	}

	return &authv1.LogoutResponse{}, nil
}
