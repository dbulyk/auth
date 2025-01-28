package model

import (
	"database/sql"
	"time"
)

type CreateUserRequest struct {
	Name            string
	Email           string
	Tag             string
	Password        string
	PasswordConfirm string
	Role            string
}

type GetUserResponse struct {
	Id        int64
	Name      string
	Email     string
	Tag       string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UpdateUserRequest struct {
	Id              int64
	Name            string
	Email           string
	Tag             string
	Password        string
	PasswordConfirm string
	Role            string
}
