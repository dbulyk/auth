package app

import (
	"context"
	"log"

	"auth/internal/api/auth"
	"auth/internal/client/db"
	"auth/internal/client/db/pg"
	"auth/internal/closer"
	"auth/internal/config"
	"auth/internal/repository"
	authRepository "auth/internal/repository/auth"
	"auth/internal/service"
	"auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	hashConfig config.HashConfig

	dbClient db.Client
	//txManager      db.TxManager
	authRepository repository.AuthRepository

	authService service.AuthService

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HashConfig() config.HashConfig {
	if s.hashConfig == nil {
		cfg, err := config.NewHashConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.hashConfig = cfg
	}

	return s.hashConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

//
//func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
//	if s.txManager == nil {
//		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
//	}
//
//	return s.txManager
//}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx), s.HashConfig().Key())
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = user.NewAuthService(s.AuthRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
