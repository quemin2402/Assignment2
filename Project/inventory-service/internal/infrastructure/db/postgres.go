package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 5
	cfg.MaxConnIdleTime = time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}
