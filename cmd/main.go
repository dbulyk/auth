package main

import (
	"auth/internal/repository"
	"auth/internal/service"
	serv "auth/internal/service/user"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"auth/internal/config"
	"auth/internal/config/env"
	desc "auth/pkg/auth_v1"
)

var configPath string

type server struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
	hashKey     string
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

	conn, err := pgx.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(pool *pgx.Conn, ctx context.Context) {
		err = pool.Close(ctx)
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn, ctx)
	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userRepo := repository.NewUserRepository()
	userService := serv.NewUserService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{userService: userService, hashKey: hashConfig.Key()})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser создает пользователя
func (s *server) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	if in.GetPassword() != in.GetPasswordConfirm() {
		return nil, status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select("email", "tag").
		From("users").
		Where(sq.Or{
			sq.Eq{"email": in.GetEmail()},
			sq.Eq{"tag": in.GetTag()}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		email string
		tag   string
	)

	err = s.db.QueryRow(ctx, query, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if len(email) > 0 {
		return nil, errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return nil, errors.New("пользователь с таким тегом уже существует")
	}

	h := hmac.New(sha256.New, []byte(s.hashKey))
	h.Write([]byte(in.GetPassword()))
	pwdHash := fmt.Sprintf("%x", h.Sum(nil))

	iBuilder := sq.Insert("users").
		Columns("name", "email", "password", "role", "tag").
		Values(in.GetName(), in.GetEmail(), pwdHash, in.GetRole(), in.GetTag()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err = iBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var userID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser получает пользователя по id
func (s *server) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	builder := sq.Select("id", "name", "email", "tag", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": in.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var (
		createdAt time.Time
		updatedAt time.Time
	)
	user := desc.GetUserResponse{}
	err = s.db.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Tag,
		&user.Role,
		&createdAt,
		&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, err
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (s *server) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	//в будущем (на 3 уроке, я так понимаю, когда будет архитектура) вынесу проверки в отдельную утилиту
	//для избежания дублирования
	if in.GetPassword() != in.GetPasswordConfirm() {
		return nil, status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select("email", "tag").
		From("users").
		Where(sq.Or{
			sq.Eq{"email": in.GetEmail()},
			sq.Eq{"tag": in.GetTag()}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		email string
		tag   string
	)

	err = s.db.QueryRow(ctx, query, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if len(email) > 0 {
		return nil, errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return nil, errors.New("пользователь с таким тегом уже существует")
	}

	h := hmac.New(sha256.New, []byte(s.hashKey))
	h.Write([]byte(in.GetPassword()))
	pwdHash := fmt.Sprintf("%x", h.Sum(nil))

	builder := sq.Update("users").
		SetMap(map[string]interface{}{
			"name":       in.GetName(),
			"email":      in.GetEmail(),
			"tag":        in.GetTag(),
			"role":       in.GetRole(),
			"updated_at": time.Now(),
			"password":   pwdHash}).
		Where(sq.Eq{"id": in.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteUser удаляет пользователя по id
func (s *server) DeleteUser(ctx context.Context, in *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	builder := sq.Delete("users").
		Where(sq.Eq{"id": in.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
