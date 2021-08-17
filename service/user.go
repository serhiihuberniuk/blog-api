package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateUser(ctx context.Context, payload models.CreateUserPayload) (string, error) {
	now := time.Now()

	user := &models.User{
		ID:        uuid.New().String(),
		Name:      payload.Name,
		Email:     payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
		Password:  payload.Password,
	}
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error occurred while hashing the password, %w", err)
	}

	user.Password = string(hashedPassword)

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}

	return user.ID, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, payload models.UpdateUserPayload) error {
	user, err := s.GetUser(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("cannot update user, %w", err)
	}

	user.Name = payload.Name
	user.Email = payload.Email
	user.UpdatedAt = time.Now()

	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
