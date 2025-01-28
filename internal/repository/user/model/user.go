package model

import (
	"database/sql"
	"time"
)

type GetUserResponse struct {
	Id        int64
	Name      string
	Email     string
	Tag       string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
