package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/serhiihuberniuk/blog-api/models"
)

type CreatePostPayload struct {
	Title       string
	Description string
	Tags        []string
}

type UpdatePostPayload struct {
	Title       string
	Description string
	Tags        []string
}

func (s *Service) CreatePost(ctx context.Context, post *models.Post, payload CreatePostPayload, user models.User) error {
	now := time.Now()
	post = &models.Post{
		ID:          uuid.New().String(),
		Title:       payload.Title,
		Description: payload.Description,
		CreatedAt:   now,
		CreatedBy:   user.ID,
		Tags:        payload.Tags,
	}

	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}
	return nil
}

func (s *Service) GetPost(ctx context.Context, post models.Post) (models.Post, error) {
	return s.repo.GetPost(ctx, post)
}

func (s *Service) UpdatePost(ctx context.Context, post *models.Post, payload UpdatePostPayload) error {
	post = &models.Post{
		Title:       payload.Title,
		Description: payload.Description,
		Tags:        payload.Tags,
	}
	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}
	return nil
}

func (s *Service) DeletePost(ctx context.Context, post models.Post) error {
	return s.repo.DeletePost(ctx, post)
}

func (s *Service) ListPosts(ctx context.Context, posts ...models.Post) ([]models.Post, error) {
	return s.repo.ListPosts(ctx, posts...)
}
