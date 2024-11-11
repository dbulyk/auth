package main

import (
	desc "auth/pkg/auth_v1"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	desc.UnimplementedAuthV1Server
	db *pgx.Conn
}

const (
	grpcPort = 50052
	dbDSN    = "host=localhost port=5432 dbname=postgres user=admin password=admin sslmode=disable"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(pool *pgx.Conn, ctx context.Context) {
		err = pool.Close(ctx)
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn, ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{db: conn})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser создает пользователя
func (s *server) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	builder := sq.Insert("users").
		Columns("name", "email", "password", "role").
		Values(in.GetName(), in.GetEmail(), in.GetPassword(), in.GetRole()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var userID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser получает пользователя по id
func (s *server) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	builder := sq.Select("name", "email", "user_tag", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": in.GetId()})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	user := &desc.GetUserResponse{}

	err = s.db.QueryRow(ctx, query, args...).Scan(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser обновляет данные пользователя
func (s *server) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	builder := sq.Update("users").
		SetMap(map[string]interface{}{
			"name":  in.GetName(),
			"email": in.GetEmail(),
			"tag":   in.GetTag(),
			"role":  in.GetRole()}).
		Where(sq.Eq{"id": in.GetId()})

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

// DeleteUser удаляет пользователя по id
func (s *server) DeleteUser(ctx context.Context, in *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	builder := sq.Delete("users").
		Where(sq.Eq{"id": in.GetId()})

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
