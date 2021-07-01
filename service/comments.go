package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
)

type CreateCommentPayload struct {
	Content  string
	PostId   string
	AuthorId string
}

type UpdateCommentPayload struct {
	CommentId string
	Content   string
}

func (s *Service) CreateComment(ctx context.Context, payload CreateCommentPayload) error {
	now := time.Now()

	comment := &models.Comment{
		ID:        uuid.New().String(),
		Content:   payload.Content,
		CreatedBy: payload.AuthorId,
		CreatedAt: now,
		PostID:    payload.PostId,
	}
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("cannot create comment: %w", err)
	}

	return fmt.Errorf("cannot create comment: %w", s.repo.CreateComment(ctx, comment))
}

func (s *Service) GetComment(ctx context.Context, commentId string) (*models.Comment, error) {
	comment, err := s.repo.GetComment(ctx, commentId)
	if err != nil {
		return nil, fmt.Errorf("cannot get comment: %w", err)
	}

	return comment, nil
}

func (s *Service) UpdateComment(ctx context.Context, payload UpdateCommentPayload) error {
	comment, err := s.GetComment(ctx, payload.CommentId)
	if err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	*comment = models.Comment{
		Content: payload.Content,
	}

	if err := comment.Validate(); err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	return fmt.Errorf("cannot update comment: %w", s.repo.UpdateComment(ctx, comment))
}

func (s *Service) DeleteComment(ctx context.Context, comment *models.Comment) error {
	return fmt.Errorf("cannot delete comment: %w", s.repo.DeleteComment(ctx, comment))
}

// todo func (s *Service) ListComments(){}
