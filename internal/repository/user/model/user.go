package model

import (
	"database/sql"
	"time"
)

// GetUserResponse является моделью данных, которая используется для получения пользователя
type GetUserResponse struct {
	ID        int64
	Name      string
	Email     string
	Tag       string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
