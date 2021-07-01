package service

import (
	"context"
	"fmt"
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
	return s.repo.CreatePost(ctx, post)
}

func (s *Service) GetPost(ctx context.Context, postId string) (*models.Post, error) {
	return s.repo.GetPost(ctx, postId)
}

func (s *Service) UpdatePost(ctx context.Context, payload UpdatePostPayload) error {
	post, err := s.GetPost(ctx, payload.PostId)
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}
	post = &models.Post{
		Title:       payload.Title,
		Description: payload.Description,
		Tags:        payload.Tags,
	}
	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}
	return s.repo.UpdatePost(ctx, post)
}

func (s *Service) DeletePost(ctx context.Context, post *models.Post) error {
	return s.repo.DeletePost(ctx, post)
}

// todo func (s *Service) ListPosts() {
