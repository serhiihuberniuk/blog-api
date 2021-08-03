package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreatePost(ctx context.Context, post *models.Post) error {
	postsCollection := usePostsCollection(r)

	if _, err := postsCollection.InsertOne(ctx, post); err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	return nil
}

func (r *Repository) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	postsCollection := usePostsCollection(r)

	var post models.Post

	if err := postsCollection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFoundPost
		}

		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return &post, nil
}

func (r *Repository) UpdatePost(ctx context.Context, post *models.Post) error {
	postsCollection := usePostsCollection(r)

	result, err := postsCollection.UpdateOne(
		ctx,
		bson.M{"_id": post.ID},
		bson.M{
			"$set": bson.M{
				"title":       post.Title,
				"description": post.Description,
				"tags":        post.Tags,
			},
		})
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	if result.MatchedCount == 0 {
		return models.ErrNotFoundPost
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, postID string) error {
	postsCollection := usePostsCollection(r)

	result, err := postsCollection.DeleteOne(ctx, bson.M{"_id": postID})
	if err != nil {
		return fmt.Errorf("cannot delete post: %w", err)
	}

	if result.DeletedCount == 0 {
		return models.ErrNotFoundPost
	}

	return nil
}

func (r *Repository) ListPosts(ctx context.Context, pagination models.Pagination,
	filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error) {
	postsCollection := usePostsCollection(r)

	filterPosts := bson.M{}
	if filter.Field != "" {
		filterPosts = bson.M{string(filter.Field): filter.Value}
	}

	opts := options.Find()
	if sort.SortByField != "" {
		opts.SetSort(bson.M{string(sort.SortByField): 1})

		if !sort.IsASC {
			opts.SetSort(bson.M{string(sort.SortByField): -1})
		}
	}

	if pagination.Offset != 0 {
		opts.SetSkip(int64(pagination.Offset))
	}

	if pagination.Limit != 0 {
		opts.SetLimit(int64(pagination.Limit))
	}

	cursor, err := postsCollection.Find(ctx, filterPosts, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot get posts: %w", err)
	}

	var posts []*models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("cannot get posts: %w", err)
	}

	return posts, nil
}
