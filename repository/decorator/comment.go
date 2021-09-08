package decorator

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
)

func (d *RepositoryCacheDecorator) CreateComment(ctx context.Context, comment *models.Comment) error {
	err := d.repository.CreateComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("error occurred in reposutory layer: %w", err)
	}

	err = d.setItemToCache(ctx, comment.ID, comment)
	if err != nil {
		return fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return nil
}
func (d *RepositoryCacheDecorator) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	var commentFromCache models.Comment
	inCache, err := d.getItemFromCache(ctx, commentID, &commentFromCache)
	if inCache {
		return &commentFromCache, nil
	}

	comment, err := d.repository.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("error occured while getting user from repository: %w", err)
	}

	err = d.setItemToCache(ctx, commentID, comment)
	if err != nil {
		return nil, fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return comment, err
}
func (d *RepositoryCacheDecorator) UpdateComment(ctx context.Context, comment *models.Comment) error {
	err := d.repository.UpdateComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.setItemToCache(ctx, comment.ID, comment)
	if err != nil {
		return fmt.Errorf("error occured while setting to cache: %w", err)
	}

	return nil
}
func (d *RepositoryCacheDecorator) DeleteComment(ctx context.Context, commentID string) error {
	err := d.repository.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.deleteItemFromCache(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error occured while deleting client from cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) ListComments(ctx context.Context, pagination models.Pagination,
	filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error) {
	return d.ListComments(ctx, pagination, filter, sort)
}
