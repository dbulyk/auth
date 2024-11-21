package model

import (
	"database/sql"
	"time"
)

// User описывает модель пользователя
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Tag       string       `db:"tag"`
	Role      int32        `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// CreateUser описывает модель для создания пользователя
type CreateUser struct {
	Name            string `db:"name"`
	Email           string `db:"email"`
	Tag             string `db:"tag"`
	Role            int32  `db:"role"`
	Password        string `db:"password"`
	PasswordConfirm string `db:"password_confirm"`
}

// UpdateUser описывает модель для обновления данных пользователя
type UpdateUser struct {
	ID              int64  `db:"id"`
	Name            string `db:"name"`
	Email           string `db:"email"`
	Tag             string `db:"tag"`
	Role            int32  `db:"role"`
	Password        string `db:"password"`
	PasswordConfirm string `db:"password_confirm"`
}
