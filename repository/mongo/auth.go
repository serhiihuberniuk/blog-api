package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) Login(ctx context.Context, email string) (*models.User, error) {
	usersCollection := useUsersCollection(r)

	var user models.User
	if err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &user, nil
}
