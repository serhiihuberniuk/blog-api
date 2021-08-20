package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (s *Service) CreateComment(ctx context.Context, payload models.CreateCommentPayload) (string, error) {
	comment := &models.Comment{
		ID:        uuid.New().String(),
		Content:   payload.Content,
		CreatedBy: s.currentUserInformationProvider.GetCurrentUserID(ctx),
		CreatedAt: time.Now(),
		PostID:    payload.PostID,
	}
	if err := comment.Validate(); err != nil {
		return "", fmt.Errorf("cannot create comment: %w", err)
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return "", fmt.Errorf("cannot create comment: %w", err)
	}

	return comment.ID, nil
}

func (s *Service) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("cannot get comment: %w", err)
	}

	return comment, nil
}

func (s *Service) UpdateComment(ctx context.Context, payload models.UpdateCommentPayload) error {
	comment, err := s.GetComment(ctx, payload.CommentID)
	if err != nil {
		return fmt.Errorf("cannot get comment: %w", err)
	}

	if !s.checkCurrentUserIsOwner(ctx, comment.CreatedBy) {
		return models.ErrNotAuthenticated
	}

	comment.Content = payload.Content

	if err := comment.Validate(); err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	if err := s.repo.UpdateComment(ctx, comment); err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	return nil
}

func (s *Service) DeleteComment(ctx context.Context, commentID string) error {
	comment, err := s.GetComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("cannot get comment: %w", err)
	}

	if !s.checkCurrentUserIsOwner(ctx, comment.CreatedBy) {
		return models.ErrNotAuthenticated
	}

	if err := s.repo.DeleteComment(ctx, commentID); err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}

	return nil
}

func (s *Service) ListComments(ctx context.Context, pagination models.Pagination,
	filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error) {
	comments, err := s.repo.ListComments(ctx, pagination, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("cannot get comments: %w", err)
	}

	return comments, nil
}
