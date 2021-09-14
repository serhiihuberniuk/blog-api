package decorator

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
)

func (d *RepositoryCacheDecorator) CreateUser(ctx context.Context, user *models.User) error {
	err := d.repository.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error occurred in reposutory layer: %w", err)
	}

	err = d.setItemToCache(ctx, user.ID, objectTypeUser, user)
	if err != nil {
		return fmt.Errorf("error occurres while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var userFromCache models.User
	if d.redisCache.Exists(ctx, objectTypeUser+userID) {
		err := d.getItemFromCache(ctx, userID, objectTypeUser, &userFromCache)
		if err != nil {
			return nil, err
		}

		return &userFromCache, nil
	}

	user, err := d.repository.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting user from repository: %w", err)
	}

	err = d.setItemToCache(ctx, userID, objectTypeUser, user)
	if err != nil {
		return nil, fmt.Errorf("error occuers while setting to cache: %w", err)
	}

	return user, err
}

func (d *RepositoryCacheDecorator) UpdateUser(ctx context.Context, user *models.User) error {
	err := d.repository.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.setItemToCache(ctx, user.ID, objectTypeUser, user)
	if err != nil {
		return fmt.Errorf("error occurred while setting to cache: %w", err)
	}

	return nil
}

func (d *RepositoryCacheDecorator) DeleteUser(ctx context.Context, userID string) error {
	err := d.repository.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error occurred in repository layer: %w", err)
	}

	err = d.deleteItemFromCache(ctx, userID, objectTypeUser)
	if err != nil {
		return fmt.Errorf("error occurred while deleting client from cache: %w", err)
	}

	return nil
}
