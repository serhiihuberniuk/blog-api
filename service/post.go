package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (s *Service) CreatePost(ctx context.Context, payload models.CreatePostPayload) error {
	post := &models.Post{
		ID:          uuid.New().String(),
		Title:       payload.Title,
		Description: payload.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   payload.AuthorID,
		Tags:        payload.Tags,
	}

	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	return nil
}

func (s *Service) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return post, nil
}

func (s *Service) UpdatePost(ctx context.Context, payload models.UpdatePostPayload) error {
	post, err := s.GetPost(ctx, payload.PostID)
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	post.Title = payload.Title
	post.Description = payload.Description
	post.Tags = payload.Tags

	if err := post.Validate(); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	return nil
}

func (s *Service) DeletePost(ctx context.Context, post *models.Post) error {
	if err := s.repo.DeletePost(ctx, post); err != nil {
		return fmt.Errorf("cannot delete post: %w", err)
	}

	return nil
}

func (s *Service) ListPosts(ctx context.Context, pagination models.Pagination,
	filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error) {
	posts, err := s.repo.ListPosts(ctx, pagination, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("cannot get posts: %w", err)
	}

	return posts, nil
}
