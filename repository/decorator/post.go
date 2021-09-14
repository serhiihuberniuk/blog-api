package decorator

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
)

func (d *RepositoryCacheDecorator) CreatePost(ctx context.Context, post *models.Post) error {
	err := d.repository.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("error occurred in reposutory layer: %w", err)
	}

	err = d.setItemToCache(ctx, post.ID, objectTypePost, post)
	if err != nil {
		return fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	var postFromCache models.Post
	if d.redisCache.Exists(ctx, objectTypePost+postID) {
		err := d.getItemFromCache(ctx, postID, objectTypePost, &postFromCache)
		if err != nil {
			return nil, fmt.Errorf("error occurred getting from cache: %w", err)
		}

		return &postFromCache, nil
	}

	post, err := d.repository.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting user from repository: %w", err)
	}

	err = d.setItemToCache(ctx, postID, objectTypePost, post)
	if err != nil {
		return nil, fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return post, err
}

func (d *RepositoryCacheDecorator) UpdatePost(ctx context.Context, post *models.Post) error {
	err := d.repository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.setItemToCache(ctx, post.ID, objectTypePost, post)
	if err != nil {
		return fmt.Errorf("error occurred while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) DeletePost(ctx context.Context, postID string) error {
	err := d.repository.DeletePost(ctx, postID)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.deleteItemFromCache(ctx, postID, objectTypePost)
	if err != nil {
		return fmt.Errorf("error occurred while deleting client from cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) ListPosts(ctx context.Context, pagination models.Pagination,
	filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error) {
	return d.repository.ListPosts(ctx, pagination, filter, sort)
}
