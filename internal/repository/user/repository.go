package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/dbulyk/platform_common/pkg/db"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
)

type repo struct {
	db db.Client
}

// NewRepository возвращает объект репозитория пользователя
func NewRepository(db db.Client) repository.UserRepository {
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

	q := db.Query{
		Name:     "user_repository.Create_CheckUser",
		QueryRaw: query,
	}

	var (
		email string
		tag   string
	)

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&email, &tag)
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

	q = db.Query{
		Name:     "user_repository.Create_InsertUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().ScanOneContext(ctx, &userID, q, args...)
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

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	user := modelRepo.GetUserResponse{}
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, in *model.UpdateUserRequest) error {
	if in.Password != in.PasswordConfirm {
		return status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select("email", "tag").
		From("users").
		Where(sq.Or{
			sq.Eq{"email": in.Email},
			sq.Eq{"tag": in.Tag}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update_CheckUser",
		QueryRaw: query,
	}

	var (
		email string
		tag   string
	)

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if len(email) > 0 {
		return errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return errors.New("пользователь с таким тегом уже существует")
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
		Where(sq.Eq{"id": in.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = builder.ToSql()
	if err != nil {
		return err
	}

	q = db.Query{
		Name:     "user_repository.Update_UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteUser(ctx context.Context, userID int64) error {
	builder := sq.Delete("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
