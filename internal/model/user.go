package model

import (
	"database/sql"
	"time"
)

// CreateUserRequest является сервисной моделью данных для запроса создания пользователя
type CreateUserRequest struct {
	Name            string
	Email           string
	Tag             string
	Password        string
	PasswordConfirm string
	Role            string
}

// GetUserResponse является сервисной моделью данных для ответа получения данных пользователя
type GetUserResponse struct {
	ID        int64
	Name      string
	Email     string
	Tag       string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UpdateUserRequest является сервисной моделью данных для запроса обновления данных пользователя
type UpdateUserRequest struct {
	ID              int64
	Name            string
	Email           string
	Tag             string
	Password        string
	PasswordConfirm string
	Role            string
}
