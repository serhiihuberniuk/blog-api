package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	Db *pgxpool.Pool
}

func NewPostgresDb(ctx context.Context, dbUrl, migrations string, version uint) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	m, err := migrate.New(migrations, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error occurred while creating migrate instances:%w", err)
	}

	if err = m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error occurred while migrating: %w", err)
	}

	return pool, nil
}

func (r *Repository) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := r.Db.Ping(ctx); err != nil {
		return fmt.Errorf("connection to database failed: %w", err)
	}

	return nil
}
