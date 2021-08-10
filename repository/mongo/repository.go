package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	Db *mongo.Database
}

func (r Repository) HealthCheck(ctx context.Context) error {
	if err := r.Db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("connection to database failed: %w", err)
	}

	return nil
}

func NewMongoDb(ctx context.Context, r *Repository) (*mongo.Database, error) {
	if err := createIndexForUsers(ctx, useUsersCollection(r)); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	if err := createIndexesForPosts(ctx, usePostsCollection(r)); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	if err := createIndexesForComments(ctx, useCommentsCollection(r)); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	return r.Db, nil
}

func useUsersCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("users")
}

func createIndexForUsers(ctx context.Context, c *mongo.Collection) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := c.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return fmt.Errorf("cannot create index, %w", err)
	}

	return nil
}

func usePostsCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("posts")
}

func createIndexesForPosts(ctx context.Context, c *mongo.Collection) error {
	mods := []mongo.IndexModel{
		{
			Keys:    bson.M{"created_by": 1},
			Options: options.Index(),
		},
		{
			Keys:    bson.M{"created_at": 1},
			Options: options.Index(),
		},
	}

	_, err := c.Indexes().CreateMany(ctx, mods)
	if err != nil {
		return fmt.Errorf("cannot create index, %w", err)
	}

	return nil
}

func useCommentsCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("comments")
}

func createIndexesForComments(ctx context.Context, c *mongo.Collection) error {
	mods := []mongo.IndexModel{
		{
			Keys:    bson.M{"created_by": 1},
			Options: options.Index(),
		},
		{
			Keys:    bson.M{"post_id": 1},
			Options: options.Index(),
		},
	}

	_, err := c.Indexes().CreateMany(ctx, mods)
	if err != nil {
		return fmt.Errorf("cannot create index, %w", err)
	}

	return nil
}
