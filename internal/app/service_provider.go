package app

import (
	"context"
	"log"

	"auth/internal/api"
	"auth/internal/client/db"
	"auth/internal/client/db/pg"
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
	dbc        db.Client

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

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepo == nil {
		r := repo.NewRepository(sp.DBClient(ctx))
		sp.userRepo = r
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
