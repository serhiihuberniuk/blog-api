package decorator

import (
	"context"
	"fmt"

	"github.com/go-redis/cache/v8"
	"github.com/serhiihuberniuk/blog-api/models"
)

type RepositoryCacheDecorator struct {
	repository  repository
	redisClient *cache.Cache
}

func NewRepositoryDecorator(r repository, c *cache.Cache) *RepositoryCacheDecorator {
	return &RepositoryCacheDecorator{
		repository:  r,
		redisClient: c,
	}
}

type repository interface {
	Login(ctx context.Context, email string) (*models.User, error)

	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, postID string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, postID string) error
	ListPosts(ctx context.Context, pagination models.Pagination,
		filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error)

	CreateComment(ctx context.Context, comment *models.Comment) error
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, commentID string) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error)
}

func (d *RepositoryCacheDecorator) setItemToCache(ctx context.Context, itemId string, value interface{}) error {
	err := d.redisClient.Set(&cache.Item{
		Ctx:   ctx,
		Key:   itemId,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) getItemFromCache(ctx context.Context, itemId string,
	destination interface{}) (bool, error) {
	if d.redisClient.Exists(ctx, itemId) {
		err := d.redisClient.Get(ctx, itemId, destination)
		if err != nil {
			return true, fmt.Errorf("error occured while getting from cache: %w", err)
		}

		return true, nil
	}

	return false, nil
}

func (d *RepositoryCacheDecorator) deleteItemFromCache(ctx context.Context, itemId string) error {
	if d.redisClient.Exists(ctx, itemId) {
		err := d.redisClient.Delete(ctx, itemId)
		if err != nil {
			return fmt.Errorf("error occurred while deleting from cache: %w", err)
		}

		return nil
	}

	return nil
}
