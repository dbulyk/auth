package model

import (
	"database/sql"
	"time"
)

// User описывает модель пользователя
type User struct {
	ID        int64
	Name      string
	Email     string
	Tag       string
	Role      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// CreateUser описывает модель для создания пользователя
type CreateUser struct {
	Name            string
	Email           string
	Tag             string
	Role            int32
	Password        string
	PasswordConfirm string
}

// UpdateUser описывает модель для обновления данных пользователя
type UpdateUser struct {
	ID              int64
	Name            string
	Email           string
	Tag             string
	Role            int32
	Password        string
	PasswordConfirm string
}
