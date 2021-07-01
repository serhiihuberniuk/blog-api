package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
)

type CreateUserPayload struct {
	Name  string
	Email string
}

type UpdateUserPayload struct {
	UserId string
	Name   string
	Email  string
}

func (s *Service) CreateUser(ctx context.Context, payload CreateUserPayload) error {
	user := &models.User{
		ID:        uuid.New().String(),
		Name:      payload.Name,
		Email:     payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, userId string) (*models.User, error) {
	return s.repo.GetUser(ctx, userId)
}

func (s *Service) UpdateUser(ctx context.Context, payload UpdateUserPayload) error {
	user, err := s.GetUser(ctx, payload.UserId)
	if err != nil {
		return fmt.Errorf("cannot update user, %w", err)
	}
	user = &models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		UpdatedAt: now,
	}
	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, user *models.User) error {
	return s.repo.DeleteUser(ctx, user)
}
