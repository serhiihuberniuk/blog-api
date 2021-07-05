package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	sql, args, err := squirrel.Insert("users").
		Values(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	sql, args, err := squirrel.Select("*").
		From("users").
		Where("id=$1", userID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	err = r.Db.QueryRow(ctx, sql, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	sql, args, err := squirrel.Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where("id=$1", user.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, user *models.User) error {
	sql, args, err := squirrel.Delete("users").
		Where("id=$1", user.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
