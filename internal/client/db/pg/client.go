package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"auth/internal/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

// New подключается к бд и возвращает клиент PG
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

// DB возвращает клиент
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close закрывает подключение к клиенту
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
