package app

import (
	"context"
	"log"

	"github.com/dbulyk/platform_common/pkg/closer"
	"github.com/dbulyk/platform_common/pkg/db"
	"github.com/dbulyk/platform_common/pkg/db/pg"
	"github.com/dbulyk/platform_common/pkg/db/trancsation"
	redigo "github.com/gomodule/redigo/redis"

	"auth/internal/api/user"
	"auth/internal/client/cache"
	"auth/internal/client/cache/redis"
	"auth/internal/config"
	"auth/internal/config/env"
	userRepo "auth/internal/repository/user"
	repo "auth/internal/repository/user/pg"
	redis2 "auth/internal/repository/user/redis"
	"auth/internal/service"
	serv "auth/internal/service/user"
)

type serviceProvider struct {
	grpcConfig  config.GRPCConfig
	pgConfig    config.PGConfig
	redisConfig config.RedisConfig

	dbc       db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepo           userRepo.Repository
	userCache          userRepo.Cache
	userService        service.UserService
	userImplementation *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		grpcConfig, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		sp.grpcConfig = grpcConfig
	}
	return sp.grpcConfig
}

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pgConfig: %v", err)
		}
		sp.pgConfig = pgConfig
	}
	return sp.pgConfig
}

func (sp *serviceProvider) RedisConfig() config.RedisConfig {
	if sp.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		sp.redisConfig = cfg
	}

	return sp.redisConfig
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbc == nil {
		conn, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = conn.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		closer.Add(func() error {
			err = conn.Close()
			if err != nil {
				return err
			}
			return nil
		})
		sp.dbc = conn
	}
	return sp.dbc
}

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = trancsation.NewTransactionManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) RedisPool() *redigo.Pool {
	if sp.redisPool == nil {
		sp.redisPool = &redigo.Pool{
			MaxIdle:     sp.RedisConfig().MaxIdle(),
			IdleTimeout: sp.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", sp.RedisConfig().Address())
			},
		}
	}

	return sp.redisPool
}

func (sp *serviceProvider) RedisClient() cache.RedisClient {
	if sp.redisClient == nil {
		sp.redisClient = redis.NewClient(sp.RedisPool(), sp.RedisConfig())
	}

	return sp.redisClient
}

func (sp *serviceProvider) UserRepository(ctx context.Context) userRepo.Repository {
	if sp.userRepo == nil {
		r := repo.NewRepository(sp.DBClient(ctx))
		sp.userRepo = r
	}
	return sp.userRepo
}

func (sp *serviceProvider) UserCache() userRepo.Cache {
	if sp.userCache == nil {
		c := redis2.NewUserCache(sp.RedisClient())
		sp.userCache = c
	}
	return sp.userCache
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		s := serv.NewUserService(sp.UserRepository(ctx), sp.TxManager(ctx), sp.UserCache())
		sp.userService = s
	}
	return sp.userService
}

func (sp *serviceProvider) UserImplementation(ctx context.Context) *user.Implementation {
	if sp.userImplementation == nil {
		i := user.NewImplementation(sp.UserService(ctx))
		sp.userImplementation = i
	}
	return sp.userImplementation
}
