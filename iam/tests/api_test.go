// Package tests содержит интеграционные тесты IAMService.
package tests

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	iampkg "github.com/Anton119/rocket-service-/iam/pkg/app"
	"github.com/Anton119/rocket-service-/iam/tests/testutil"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	commonv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/common/v1"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

const (
	bufSize        = 1024 * 1024
	testPassword   = "password123"
	testSessionTTL = time.Hour
)

var (
	infra      *iampkg.Infra
	authClient authv1.AuthServiceClient
	userClient userv1.UserServiceClient
)

func initTestEnv() {
	if os.Getenv("IAM_DB_URI") == "" && os.Getenv("DB_URI") == "" {
		os.Setenv("IAM_DB_URI", iampkg.IAMDBURI())
	}
	if os.Getenv("REDIS_HOST") == "" {
		os.Setenv("REDIS_HOST", "localhost")
	}
	if os.Getenv("REDIS_PORT") == "" {
		os.Setenv("REDIS_PORT", "6379")
	}
}

func TestMain(m *testing.M) {
	initTestEnv()

	ctx := context.Background()

	var err error
	infra, err = iampkg.OpenInfra(ctx)
	if err != nil {
		panic(err)
	}

	if err = infra.Reset(ctx); err != nil {
		panic(err)
	}

	lis := bufconn.Listen(bufSize)
	grpcServer, err := iampkg.NewGRPCServer(infra, testSessionTTL)
	if err != nil {
		panic(err)
	}

	go func() {
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			panic(fmt.Sprintf("gRPC server: %v", serveErr))
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	authClient = authv1.NewAuthServiceClient(conn)
	userClient = userv1.NewUserServiceClient(conn)

	code := m.Run()

	conn.Close()
	grpcServer.Stop()
	infra.Close()
	os.Exit(code)
}

func resetData(t *testing.T) {
	t.Helper()
	require.NoError(t, infra.Reset(context.Background()))
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

func registerUser(t *testing.T, login string) string {
	t.Helper()

	resp, err := userClient.Register(context.Background(), registerRequest(login, testPassword))
	require.NoError(t, err)

	return resp.GetUserUuid()
}

func loginUser(t *testing.T, login string) string {
	t.Helper()

	resp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Login:    login,
		Password: testPassword,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetSessionUuid())

	return resp.GetSessionUuid()
}

func uniqueLogin(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, uuid.NewString()[:8])
}

func TestRegister_Success(t *testing.T) {
	resetData(t)

	login := uniqueLogin("register")
	resp, err := userClient.Register(context.Background(), registerRequest(login, testPassword))
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUserUuid())
}

func TestRegister_DuplicateLogin(t *testing.T) {
	resetData(t)

	login := uniqueLogin("duplicate")
	_, err := userClient.Register(context.Background(), registerRequest(login, testPassword))
	require.NoError(t, err)

	_, err = userClient.Register(context.Background(), registerRequest(login, testPassword))
	testutil.AssertGRPCStatus(t, err, codes.AlreadyExists)
}

func TestRegister_WeakPassword(t *testing.T) {
	resetData(t)

	_, err := userClient.Register(context.Background(), registerRequest(uniqueLogin("weak"), "short"))
	testutil.AssertGRPCStatus(t, err, codes.InvalidArgument)
}

func TestLogin_Success(t *testing.T) {
	resetData(t)

	login := uniqueLogin("login")
	registerUser(t, login)

	resp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Login:    login,
		Password: testPassword,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetSessionUuid())
}

func TestLogin_WrongPassword(t *testing.T) {
	resetData(t)

	login := uniqueLogin("wrong-pass")
	registerUser(t, login)

	_, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Login:    login,
		Password: "wrongpassword",
	})
	testutil.AssertGRPCStatus(t, err, codes.Unauthenticated)
}

func TestWhoami_Success(t *testing.T) {
	resetData(t)

	login := uniqueLogin("whoami")
	registerUser(t, login)
	sessionUUID := loginUser(t, login)

	resp, err := authClient.Whoami(context.Background(), &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	require.NoError(t, err)
	assert.Equal(t, sessionUUID, resp.GetSession().GetUuid())
	assert.Equal(t, login, resp.GetUser().GetInfo().GetLogin())
}

func TestWhoami_InvalidSession(t *testing.T) {
	resetData(t)

	_, err := authClient.Whoami(context.Background(), &authv1.WhoamiRequest{
		SessionUuid: uuid.NewString(),
	})
	testutil.AssertGRPCStatus(t, err, codes.Unauthenticated)
}

func TestLogout_Success(t *testing.T) {
	resetData(t)

	login := uniqueLogin("logout")
	registerUser(t, login)
	sessionUUID := loginUser(t, login)

	_, err := authClient.Logout(context.Background(), &authv1.LogoutRequest{
		SessionUuid: sessionUUID,
	})
	require.NoError(t, err)

	_, err = authClient.Whoami(context.Background(), &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	testutil.AssertGRPCStatus(t, err, codes.Unauthenticated)
}

func TestLogout_Idempotent(t *testing.T) {
	resetData(t)

	login := uniqueLogin("logout-idem")
	registerUser(t, login)
	sessionUUID := loginUser(t, login)

	_, err := authClient.Logout(context.Background(), &authv1.LogoutRequest{
		SessionUuid: sessionUUID,
	})
	require.NoError(t, err)

	_, err = authClient.Logout(context.Background(), &authv1.LogoutRequest{
		SessionUuid: sessionUUID,
	})
	require.NoError(t, err)
}

func TestGetUser_Success(t *testing.T) {
	resetData(t)

	login := uniqueLogin("get-user")
	userUUID := registerUser(t, login)

	resp, err := userClient.GetUser(context.Background(), &userv1.GetUserRequest{
		UserUuid: userUUID,
	})
	require.NoError(t, err)
	assert.Equal(t, userUUID, resp.GetUser().GetUuid())
	assert.Equal(t, login, resp.GetUser().GetInfo().GetLogin())
}

func TestFullFlow(t *testing.T) {
	resetData(t)
	ctx := context.Background()
	login := uniqueLogin("fullflow")

	userUUID := registerUser(t, login)
	sessionUUID := loginUser(t, login)

	whoamiResp, err := authClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	require.NoError(t, err)
	assert.Equal(t, login, whoamiResp.GetUser().GetInfo().GetLogin())

	getUserResp, err := userClient.GetUser(ctx, &userv1.GetUserRequest{UserUuid: userUUID})
	require.NoError(t, err)
	assert.Equal(t, userUUID, getUserResp.GetUser().GetUuid())

	_, err = authClient.Logout(ctx, &authv1.LogoutRequest{SessionUuid: sessionUUID})
	require.NoError(t, err)

	_, err = authClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	testutil.AssertGRPCStatus(t, err, codes.Unauthenticated)
}
