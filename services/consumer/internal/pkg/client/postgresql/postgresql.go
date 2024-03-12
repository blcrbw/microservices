package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"consumer/internal/config"
	"consumer/internal/utils"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.DbConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?pool_max_conns=%s&pool_min_conns=%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database, sc.PoolMaxConns, sc.PoolMinConns)
	err = reuse.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		panic("error do with tries postgresql")
	}

	return pool, nil
}
