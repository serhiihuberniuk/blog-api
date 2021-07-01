package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/serhiihuberniuk/blog-api/models"
)

type CreateUserPayload struct {
	name  string
	email string
}

type UpdateUserPayload struct {
	name  string
	email string
}

func (s *Service) CreateUser(ctx context.Context, user *models.User, payload CreateUserPayload) error {
	now := time.Now()
	user = &models.User{
		ID:        uuid.New().String(),
		Name:      payload.name,
		Email:     payload.email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, user models.User) (models.User, error) {
	return s.repo.GetUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user *models.User, payload UpdateUserPayload) error {
	now := time.Now()
	user = &models.User{
		Name:      payload.name,
		Email:     payload.email,
		UpdatedAt: now,
	}
	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, user models.User) error {
	return s.repo.DeleteUser(ctx, user)
}
