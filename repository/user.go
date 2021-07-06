package repository

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	sql := "INSERT INTO users VALUES ($1, $2, $3, $4, $5)"

	_, err := r.Db.Exec(ctx, sql, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	sql := "SELECT * FROM users WHERE id=$1"

	err := pgxscan.Get(ctx, r.Db, &user, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	sql := "UPDATE users SET name=$1, email=$2, updated_at=$3 WHERE id=$4"

	_, err := r.Db.Exec(ctx, sql, user.Name, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, user *models.User) error {
	sql := "DELETE FROM users WHERE id=$1"

	_, err := r.Db.Exec(ctx, sql, user.ID)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
