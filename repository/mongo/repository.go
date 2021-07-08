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

func NewMongoDb(ctx context.Context, db *mongo.Database) (*mongo.Database, error) {
	usersCollection := db.Collection("users")
	if err := createIndexForUsers(ctx, usersCollection); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	postsCollection := db.Collection("posts")
	if err := createIndexesForPosts(ctx, postsCollection); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	commentsCollection := db.Collection("comments")
	if err := createIndexesForComments(ctx, commentsCollection); err != nil {
		return nil, fmt.Errorf("cannot create database: %w", err)
	}

	return db, nil
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
	mod := mongo.IndexModel{
		Keys: bson.M{
			"created_by": 1,
			"created_at": 1,
		},
		Options: options.Index(),
	}

	_, err := c.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return fmt.Errorf("cannot create index, %w", err)
	}

	return nil
}

func useCommentsCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("comments")
}

func createIndexesForComments(ctx context.Context, c *mongo.Collection) error {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"post_id":    1,
			"created_by": 1,
		},
		Options: options.Index(),
	}

	_, err := c.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return fmt.Errorf("cannot create index, %w", err)
	}

	return nil
}
