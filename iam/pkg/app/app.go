// Package app содержит вспомогательные функции сборки IAMService для интеграционных тестов.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	authapi "github.com/Anton119/rocket-service-/iam/internal/api/auth/v1"
	userapi "github.com/Anton119/rocket-service-/iam/internal/api/user/v1"
	iaminterceptor "github.com/Anton119/rocket-service-/iam/internal/interceptor"
	sessionrepo "github.com/Anton119/rocket-service-/iam/internal/repository/session"
	userrepo "github.com/Anton119/rocket-service-/iam/internal/repository/user"
	iamsvc "github.com/Anton119/rocket-service-/iam/internal/service/iam"
	redisclient "github.com/Anton119/rocket-service-/platform/pkg/redis"
	"github.com/Anton119/rocket-service-/shared/pkg/grpc/interceptor"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

const defaultSessionTTL = 24 * time.Hour

// Infra — PostgreSQL и Redis для IAM.
type Infra struct {
	Pool  *pgxpool.Pool
	Redis *redis.Client
}

// IAMDBURI возвращает DSN PostgreSQL IAM (IAM_DB_URI, DB_URI или дефолт для локальных тестов).
func IAMDBURI() string {
	if dsn := os.Getenv("IAM_DB_URI"); dsn != "" {
		return dsn
	}

	if dsn := os.Getenv("DB_URI"); dsn != "" {
		return dsn
	}

	return "postgres://iam-service-user:iam-service-password@localhost:5434/iam-service?sslmode=disable"
}

// RedisAddr возвращает адрес Redis (host:port).
func RedisAddr() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return host + ":" + port
}

// OpenInfra подключается к PostgreSQL и Redis.
func OpenInfra(ctx context.Context) (*Infra, error) {
	pool, err := pgxpool.New(ctx, IAMDBURI())
	if err != nil {
		return nil, fmt.Errorf("создание пула PostgreSQL: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("проверка PostgreSQL: %w", err)
	}

	rdb, err := redisclient.NewClient(&redis.Options{
		Addr: RedisAddr(),
	}, slog.Default())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("подключение к Redis: %w", err)
	}

	return &Infra{Pool: pool, Redis: rdb}, nil
}

// Close закрывает соединения с PostgreSQL и Redis.
func (i *Infra) Close() {
	if i == nil {
		return
	}

	if i.Redis != nil {
		_ = i.Redis.Close() //nolint:gosec // cleanup в тестах, ошибку Close игнорируем.
	}

	if i.Pool != nil {
		i.Pool.Close()
	}
}

// Reset очищает данные IAM перед/между тестами.
func (i *Infra) Reset(ctx context.Context) error {
	if _, err := i.Pool.Exec(ctx, `TRUNCATE users`); err != nil {
		return fmt.Errorf("очистить users: %w", err)
	}

	if err := i.Redis.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf("очистить Redis: %w", err)
	}

	return nil
}

// NewGRPCServer собирает gRPC-сервер IAM с реальными репозиториями.
func NewGRPCServer(infra *Infra, sessionTTL time.Duration) (*grpc.Server, error) {
	if sessionTTL == 0 {
		sessionTTL = defaultSessionTTL
	}

	pvUnary, err := interceptor.UnaryProtovalidateInterceptor()
	if err != nil {
		return nil, fmt.Errorf("protovalidate interceptor: %w", err)
	}

	svc := iamsvc.NewService(
		userrepo.New(infra.Pool),
		sessionrepo.New(infra.Redis),
		sessionTTL,
	)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.RecoveryInterceptor(),
		pvUnary,
		interceptor.LoggerInterceptor(),
		iaminterceptor.ErrorInterceptor(),
	))

	authv1.RegisterAuthServiceServer(srv, authapi.NewAPI(svc))
	userv1.RegisterUserServiceServer(srv, userapi.NewAPI(svc))

	return srv, nil
}
