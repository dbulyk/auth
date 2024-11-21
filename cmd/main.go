package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"auth/internal/config"
	"auth/internal/config/env"
	"auth/internal/repository/auth"
	"auth/internal/repository/auth/model"
	desc "auth/pkg/auth_v1"
)

var configPath string

type server struct {
	desc.UnimplementedAuthV1Server
	repo *auth.Repo
}

// init записывает параметр конфига
func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	hashConfig, err := env.NewHashConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authRepo := auth.NewRepository(pool, hashConfig.Key())
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{repo: authRepo})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser создает пользователя
func (s *server) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	user := model.CreateUser{
		Name:            in.GetName(),
		Email:           in.GetEmail(),
		Tag:             in.GetTag(),
		Role:            int32(in.GetRole()),
		Password:        in.GetPassword(),
		PasswordConfirm: in.GetPasswordConfirm(),
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания пользователя: %v", err)
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser получает пользователя по id
func (s *server) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := s.repo.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных пользователя: %v", err)
	}

	createdAt := timestamppb.New(user.CreatedAt)
	updatedAt := timestamppb.New(user.UpdatedAt.Time)
	res := &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Tag:       user.Tag,
		Role:      desc.Role(user.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return res, nil
}

// UpdateUser обновляет данные пользователя
func (s *server) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	updateUser := model.UpdateUser{
		Name:            in.GetName(),
		Email:           in.GetEmail(),
		Tag:             in.GetTag(),
		Password:        in.GetPassword(),
		PasswordConfirm: in.GetPasswordConfirm(),
	}
	err := s.repo.UpdateUser(ctx, updateUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления данных пользователя: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser удаляет пользователя по id
func (s *server) DeleteUser(ctx context.Context, in *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.repo.DeleteUser(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления пользователя: %v", err)
	}

	return &emptypb.Empty{}, nil
}
