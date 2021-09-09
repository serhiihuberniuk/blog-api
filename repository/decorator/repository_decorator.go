package decorator

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/serhiihuberniuk/blog-api/models"
)

type RepositoryCacheDecorator struct {
	repository  repository
	redisClient *redis.Client
	redisCache  *cache.Cache
}

func NewRepositoryCacheDecorator(r repository, redisAddress string) *RepositoryCacheDecorator {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		DB:       0,
		Password: "",
	})

	newCache := cache.New(&cache.Options{
		Redis: client,
	})

	return &RepositoryCacheDecorator{
		repository:  r,
		redisClient: client,
		redisCache:  newCache,
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

func (d *RepositoryCacheDecorator) setItemToCache(ctx context.Context, itemID string, value interface{}) error {
	err := d.redisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   itemID,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) getItemFromCache(ctx context.Context, itemID string, destination interface{}) error {
	if err := d.redisCache.Get(ctx, itemID, destination); err != nil {
		return fmt.Errorf("error occured while getting from cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) deleteItemFromCache(ctx context.Context, itemID string) error {
	if d.redisCache.Exists(ctx, itemID) {
		err := d.redisCache.Delete(ctx, itemID)
		if err != nil {
			return fmt.Errorf("error occurred while deleting from cache: %w", err)
		}

		return nil
	}

	return nil
}

func (d *RepositoryCacheDecorator) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := d.redisClient.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("connection to redis db failed: %w", err)
	}

	return nil
}
