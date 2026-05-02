package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("config inválida de DB: %w", err)
    }

    cfg.MaxConns = 10
    cfg.MinConns = 2

    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, fmt.Errorf("no se pudo conectar a postgres: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping fallido: %w", err)
    }

    return pool, nil
}