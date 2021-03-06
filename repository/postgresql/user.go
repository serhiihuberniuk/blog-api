package repository

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	const sql = "INSERT INTO users (id, name, email, created_at, updated_at, password) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.Db.Exec(ctx, sql, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.Password)
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	const sql = "SELECT id, name, email, created_at, updated_at FROM users WHERE id=$1"

	var user models.User

	err := pgxscan.Get(ctx, r.Db, &user, sql, userID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	const sql = "UPDATE users SET name=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5"

	result, err := r.Db.Exec(ctx, sql, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID string) error {
	const sql = "DELETE FROM users WHERE id=$1"

	result, err := r.Db.Exec(ctx, sql, userID)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
