package repository

import (
	"context"
	"fmt"
	"github.com/serhiihuberniuk/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
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
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	usersCollection := useUsersCollection(r)

	if _, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"name":       user.Name,
				"email":      user.Email,
				"updated_at": user.UpdatedAt,
			},
		}); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID string) error {
	usersCollection := useUsersCollection(r)

	if _, err := usersCollection.DeleteOne(ctx, bson.M{"_id": userID}); err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
