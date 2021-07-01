package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
)

type CreatePostPayload struct {
	Title       string
	Description string
	Tags        []string
	AuthorId    string
}

type UpdatePostPayload struct {
	PostId      string
	Title       string
	Description string
	Tags        []string
}

func (s *Service) CreatePost(ctx context.Context, payload CreatePostPayload) error {
	now := time.Now()

	post := &models.Post{
		ID:          uuid.New().String(),
		Title:       payload.Title,
		Description: payload.Description,
		CreatedAt:   now,
		CreatedBy:   payload.AuthorId,
		Tags:        payload.Tags,
	}

	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	return fmt.Errorf("cannot create post: %w", s.repo.CreatePost(ctx, post))
}

func (s *Service) GetPost(ctx context.Context, postId string) (*models.Post, error) {
	post, err := s.repo.GetPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return post, nil
}

func (s *Service) UpdatePost(ctx context.Context, payload UpdatePostPayload) error {
	post, err := s.GetPost(ctx, payload.PostId)
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	*post = models.Post{
		Title:       payload.Title,
		Description: payload.Description,
		Tags:        payload.Tags,
	}

	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	return fmt.Errorf("cannot update post: %w", s.repo.UpdatePost(ctx, post))
}

func (s *Service) DeletePost(ctx context.Context, post *models.Post) error {
	return fmt.Errorf("cannot delete post: %w", s.repo.DeletePost(ctx, post))
}

// todo func (s *Service) ListPosts() {
