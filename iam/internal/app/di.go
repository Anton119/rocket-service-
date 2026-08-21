package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	authapi "github.com/Anton119/rocket-service-/iam/internal/api/auth/v1"
	userapi "github.com/Anton119/rocket-service-/iam/internal/api/user/v1"
	"github.com/Anton119/rocket-service-/iam/internal/config"
	sessionrepo "github.com/Anton119/rocket-service-/iam/internal/repository/session"
	userrepo "github.com/Anton119/rocket-service-/iam/internal/repository/user"
	iamsvc "github.com/Anton119/rocket-service-/iam/internal/service/iam"
	"github.com/Anton119/rocket-service-/platform/pkg/closer"
	redisclient "github.com/Anton119/rocket-service-/platform/pkg/redis"
	authv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/auth/v1"
	userv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/user/v1"
)

func newGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	all := append([]grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}, opts...)

	return grpc.NewServer(all...)
}

type diContainer struct {
	pgPool  *pgxpool.Pool
	redis   *redis.Client
	iamSvc  *iamsvc.Service
	authAPI authv1.AuthServiceServer
	userAPI userv1.UserServiceServer
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		dsn, err := config.AppConfig().PG.DSN()
		if err != nil {
			slog.Error("получение DSN", "error", err)
			os.Exit(1)
		}

		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("создание пула соединений", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			slog.Error("проверка соединения с БД", "error", err)
			os.Exit(1)
		}
		slog.Info("подключение к PostgreSQL установлено")

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}
	return d.pgPool
}

func (d *diContainer) RedisClient(ctx context.Context) *redis.Client {
	if d.redis == nil {
		cfg := config.AppConfig().Redis
		client, err := redisclient.NewClient(&redis.Options{
			Addr: cfg.Address(),
		}, slog.Default())
		if err != nil {
			slog.Error("создание Redis-клиента", "error", err)
			os.Exit(1)
		}

		closer.Add("Redis client", func(_ context.Context) error {
			return client.Close()
		})

		d.redis = client
	}

	return d.redis
}

func (d *diContainer) IAMService(ctx context.Context) *iamsvc.Service {
	if d.iamSvc == nil {
		d.iamSvc = iamsvc.NewService(
			userrepo.New(d.PGPool(ctx)),
			sessionrepo.New(d.RedisClient(ctx)),
			config.AppConfig().Session.TTL,
		)
	}
	return d.iamSvc
}

func (d *diContainer) AuthAPI(ctx context.Context) authv1.AuthServiceServer {
	if d.authAPI == nil {
		d.authAPI = authapi.NewAPI(d.IAMService(ctx))
	}
	return d.authAPI
}

func (d *diContainer) UserAPI(ctx context.Context) userv1.UserServiceServer {
	if d.userAPI == nil {
		d.userAPI = userapi.NewAPI(d.IAMService(ctx))
	}
	return d.userAPI
}
