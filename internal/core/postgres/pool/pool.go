package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	core_postgres "github.com/Fitray/sentinel-service/internal/core/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
	GetTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	Timeout time.Duration
}

func NewConnectionPool(config core_postgres.Config) (ConnectionPool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Timeout,
	)
	defer cancel()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DB,
	)

	fmt.Println(connString)

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return ConnectionPool{},
			fmt.Errorf("failed to parse postgre connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return ConnectionPool{},
			fmt.Errorf("failed to connect to pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return ConnectionPool{},
			fmt.Errorf("failed to ping pool: %w", err)
	}

	return ConnectionPool{
		Pool:    pool,
		Timeout: config.Timeout,
	}, nil
}

func (c ConnectionPool) GetTimeout() time.Duration {
	return c.Timeout
}
