package pg

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/client/db"
)

type key string

const (
	// TxKey константа для транзакции
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB создаёт объект бд
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext является обёрткой для получения одной записи из бд и парсинга её в структуру
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	//logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext является оберткой для получения нескольких записей из бд и парсинга их в структуру
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	//logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext является оберткой для выполнения запроса к бд
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	//logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext является оберткой для запроса к бд в транзакционном режиме с возвратом множественного результата
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	//logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext является оберткой для запроса к бд в транзакционном режиме с возвратом одиночного результата
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	//logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// BeginTx открывает транзакцию
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping проверяет соединение с бд
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close закрывает соединение с бд
func (p *pg) Close() {
	p.dbc.Close()
}

//func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
//	return context.WithValue(ctx, TxKey, tx)
//}

//func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
//	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
//	log.Println(
//		ctx,
//		fmt.Sprintf("sql: %s", q.Name),
//		fmt.Sprintf("query: %s", prettyQuery),
//	)
//}
