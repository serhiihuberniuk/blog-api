package repository

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) Login(ctx context.Context, email, password string) (*models.User, error) {
	const sql = "SELECT id, name, email, created_at, updated_at, password FROM users WHERE email=$1"

	var user models.User

	err := pgxscan.Get(ctx, r.Db, &user, sql, email)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("authentication failed, %w", err)
	}

	return &user, nil
}
