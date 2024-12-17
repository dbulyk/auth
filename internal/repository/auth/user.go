package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository"
)

var _ repository.User = (*repoUser)(nil)

type repoUser struct {
	db      db.Client
	hashKey string
}

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	tagColumn       = "tag"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

func (r repoUser) CreateUser(ctx context.Context, user model.CreateUser) (id int64, err error) {
	if user.Password != user.PasswordConfirm {
		return -1, status.Error(codes.FailedPrecondition, "пароли не совпадают")
	}

	sBuilder := sq.Select(emailColumn, tagColumn).
		From(tableName).
		Where(sq.Or{
			sq.Eq{emailColumn: user.Email},
			sq.Eq{tagColumn: user.Tag}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()
	if err != nil {
		return -1, err
	}

	var (
		email string
		tag   string
	)

	q := db.Query{
		Name:     "auth_repository.CheckUser",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&email, &tag)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return -1, err
	}

	if len(email) > 0 {
		return -1, errors.New("пользователь с таким email уже существует")
	} else if len(tag) > 0 {
		return -1, errors.New("пользователь с таким тегом уже существует")
	}

	h := hmac.New(sha256.New, []byte(r.hashKey))
	h.Write([]byte(user.Password))
	pwdHash := fmt.Sprintf("%x", h.Sum(nil))

	iBuilder := sq.Insert(tableName).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, tagColumn).
		Values(user.Name, user.Email, pwdHash, user.Role, user.Tag).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err = iBuilder.ToSql()
	if err != nil {
		return -1, err
	}

	q = db.Query{
		Name:     "auth_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, status.Errorf(codes.NotFound, "пользователь не создан")
		}
		return -1, err
	}
	return userID, nil
}

func (r repoUser) UpdateUser(ctx context.Context, user model.UpdateUser) error {
	var pwdHash string
	if len(user.Password) > 0 {
		if user.Password != user.PasswordConfirm {
			return status.Error(codes.FailedPrecondition, "пароли не совпадают")
		}

		sBuilder := sq.Select(emailColumn, tagColumn).
			From(tableName).
			Where(sq.Or{
				sq.Eq{emailColumn: user.Email},
				sq.Eq{tagColumn: user.Tag}}).
			PlaceholderFormat(sq.Dollar)

		query, args, err := sBuilder.ToSql()
		if err != nil {
			return err
		}

		var (
			email string
			tag   string
		)

		q := db.Query{
			Name:     "auth_repository.GetUser",
			QueryRaw: query,
		}

		err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&email, &tag)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		if len(email) > 0 {
			return errors.New("пользователь с таким email уже существует")
		} else if len(tag) > 0 {
			return errors.New("пользователь с таким тегом уже существует")
		}

		h := hmac.New(sha256.New, []byte(r.hashKey))
		h.Write([]byte(user.Password))
		pwdHash = fmt.Sprintf("%x", h.Sum(nil))
	}

	var m map[string]interface{}

	if len(pwdHash) > 0 {
		m[passwordColumn] = pwdHash
	}
	if len(user.Name) > 0 {
		m[nameColumn] = user.Name
	}
	if len(user.Email) > 0 {
		m[emailColumn] = user.Email
	}
	if len(user.Tag) > 0 {
		m[tagColumn] = user.Tag
	}
	if user.Role > 0 {
		m[roleColumn] = user.Role
	}

	builder := sq.Update(tableName).
		SetMap(m).
		Where(sq.Eq{idColumn: user.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r repoUser) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, tagColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.GetUser",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Tag, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser удаляет пользователя по его id
func (r repoUser) DeleteUser(ctx context.Context, id int64) (err error) {
	builder := sq.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
