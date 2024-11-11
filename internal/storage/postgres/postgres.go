package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func NewDataBase(ctx context.Context, url string) (*sqlx.DB, error) {
	connPoolConf, err := pgx.ParseConnectionString(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection pool config: %w", err)
	}

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connPoolConf,
		MaxConnections: 20,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection pool: %w", err)
	}

	db := sqlx.NewDb(stdlib.OpenDBFromPool(
		connPool,
		stdlib.OptionPreferSimpleProtocol(true),
	), "pgx")

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	db.Mapper = reflectx.NewMapperFunc("db", strings.ToLower)

	return db, nil
}
