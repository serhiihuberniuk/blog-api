package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type Repository struct {
	Db *pgxpool.Pool
}

func NewPostgresDb(ctx context.Context, dbUrl, init string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	sql, err := os.ReadFile(init)
	if err != nil {
		return nil, fmt.Errorf("cannot read postgres db initialisation file: %w", err)
	}

	_, err = pool.Exec(ctx, string(sql))
	if err != nil {
		return nil, fmt.Errorf("cannot initiate postgres db: %w", err)
	}

	return pool, nil
}

func (r *Repository) HealthCheck(ctx context.Context) error {
	if err := r.Db.Ping(ctx); err != nil {
		return fmt.Errorf("connection to database failed: %w", err)
	}

	return nil
}
