package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/serhiihuberniuk/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	usersCollection := useUsersCollection(r)

	if _, err := usersCollection.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	usersCollection := useUsersCollection(r)

	var user models.User
	if err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	usersCollection := useUsersCollection(r)

	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"name":       user.Name,
				"email":      user.Email,
				"updated_at": user.UpdatedAt,
			},
		})
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID string) error {
	usersCollection := useUsersCollection(r)

	result, err := usersCollection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return models.ErrNotFound
	}

	return nil
}
