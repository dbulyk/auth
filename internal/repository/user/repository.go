package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, in *model.CreateUserRequest) (int64, error) {
	if in.Password != in.PasswordConfirm {
		return 0, status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select("email", "tag").
		From("users").
		Where(sq.Or{
			sq.Eq{"email": in.Email},
			sq.Eq{"tag": in.Tag}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var (
		email string
		tag   string
	)

	err = r.db.QueryRow(ctx, query, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	if len(email) > 0 {
		return 0, errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return 0, errors.New("пользователь с таким тегом уже существует")
	}

	h := hmac.New(sha256.New, []byte("test"))
	h.Write([]byte(in.Password))
	pwdHash := fmt.Sprintf("%x", h.Sum(nil))

	iBuilder := sq.Insert("users").
		Columns("name", "email", "password", "role", "tag").
		Values(in.Name, in.Email, pwdHash, in.Role, in.Tag).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err = iBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return 0, err
	}

	return userID, nil
}

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.GetUserResponse, error) {
	builder := sq.Select("id", "name", "email", "tag", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var (
		createdAt time.Time
		updatedAt time.Time
	)
	user := modelRepo.GetUserResponse{}
	err = r.db.QueryRow(ctx, query, args...).Scan(
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

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error) {
	if in.Password != in.PasswordConfirm {
		return nil, status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select("email", "tag").
		From("users").
		Where(sq.Or{
			sq.Eq{"email": in.Email},
			sq.Eq{"tag": in.Tag}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		email string
		tag   string
	)

	err = r.db.QueryRow(ctx, query, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if len(email) > 0 {
		return nil, errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return nil, errors.New("пользователь с таким тегом уже существует")
	}

	h := hmac.New(sha256.New, []byte("test")) //TODO добавить обработку хеша
	h.Write([]byte(in.Password))
	pwdHash := fmt.Sprintf("%x", h.Sum(nil))

	builder := sq.Update("users").
		SetMap(map[string]interface{}{
			"name":       in.Name,
			"email":      in.Email,
			"tag":        in.Tag,
			"role":       in.Role,
			"updated_at": time.Now(),
			"password":   pwdHash}).
		Where(sq.Eq{"id": in.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	builder := sq.Delete("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
