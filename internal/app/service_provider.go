package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/api"
	"auth/internal/closer"
	"auth/internal/config"
	"auth/internal/config/env"
	"auth/internal/repository"
	repo "auth/internal/repository/user"
	"auth/internal/service"
	serv "auth/internal/service/user"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig
	pgConfig   config.PGConfig
	pgPool     *pgxpool.Pool

	userRepo           repository.UserRepository
	userService        service.UserService
	userImplementation *api.Implementation
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

func (sp *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if sp.pgPool == nil {
		conn, err := pgxpool.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		closer.Add(func() error {
			conn.Close()
			return nil
		})
		sp.pgPool = conn
	}
	return sp.pgPool
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepo == nil {
		repo := repo.NewRepository(sp.PGPool(ctx))
		sp.userRepo = repo
	}
	return sp.userRepo
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		s := serv.NewUserService(sp.UserRepository(ctx))
		sp.userService = s
	}
	return sp.userService
}

func (sp *serviceProvider) UserImplementation(ctx context.Context) *api.Implementation {
	if sp.userImplementation == nil {
		i := api.NewImplementation(sp.UserService(ctx))
		sp.userImplementation = i
	}
	return sp.userImplementation
}
