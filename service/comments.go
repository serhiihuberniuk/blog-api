package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
	"time"
)

type CreateCommentPayload struct {
	content string
}

type UpdateCommentPayload struct {
	content string
}

func (s *Service) CreateComment(ctx context.Context, payload CreateCommentPayload, comment *models.Comment, user models.User, post models.Post) error {
	now := time.Now()
	comment = &models.Comment{
		ID:        uuid.New().String(),
		Content:   payload.content,
		CreatedBy: user.ID,
		CreatedAt: now,
		PostID:    post.ID,
	}
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("cannot create comment: %w", err)
	}
	return nil
}

func (s *Service) GetComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	return s.repo.GetComment(ctx, comment)
}

func (s *Service) UpdateComment(ctx context.Context, comment *models.Comment, payload UpdateCommentPayload) error {
	comment = &models.Comment{
		Content: payload.content,
	}
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}
	return nil
}

func (s *Service) DeleteComment(ctx context.Context, comment models.Comment) error {
	return s.repo.DeleteComment(ctx, comment)
}

func (s *Service) ListComments(ctx context.Context, comments ...models.Comment) ([]models.Comment, error) {
	return s.repo.ListComments(ctx, comments...)
}
