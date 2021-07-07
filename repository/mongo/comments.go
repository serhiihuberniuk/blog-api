package repository

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateComment(ctx context.Context, comment *models.Comment) error {
	commentsCollection := useCommentsCollection(r)

	if _, err := commentsCollection.InsertOne(ctx, comment); err != nil {
		return fmt.Errorf("cannot create comment: %w", err)
	}

	return nil
}

func (r *Repository) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	commentsCollection := useCommentsCollection(r)

	var comment models.Comment

	if err := commentsCollection.FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment); err != nil {
		return nil, fmt.Errorf("cannot get comment: %w", err)
	}

	return &comment, nil
}

func (r *Repository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	commentsCollection := useCommentsCollection(r)

	if _, err := commentsCollection.UpdateOne(
		ctx,
		bson.M{"_id": comment.ID},
		bson.M{
			"$set": bson.M{
				"content": comment.Content,
			},
		}); err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, comment *models.Comment) error {
	commentsCollection := useCommentsCollection(r)

	if _, err := commentsCollection.DeleteOne(ctx, bson.M{"_id": comment.ID}); err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}

	return nil
}

func (r *Repository) ListComments(ctx context.Context, pagination models.Pagination,
	filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error) {
	commentsCollection := useCommentsCollection(r)

	filterComments := bson.M{}
	if filter.Field != "" {
		filterComments = bson.M{string(filter.Field): filter.Value}
	}

	opts := options.Find()
	if sort.Field != "" {
		opts.SetSort(bson.M{string(sort.Field): 1})

		if !sort.IsASC {
			opts.SetSort(bson.M{string(sort.Field): -1})
		}
	}

	if pagination.Offset != 0 {
		opts.SetSkip(int64(pagination.Offset))
	}

	if pagination.Limit != 0 {
		opts.SetLimit(int64(pagination.Limit))
	}

	cursor, err := commentsCollection.Find(ctx, filterComments, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot get comments: %w", err)
	}

	var comments []*models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("cannot get comments: %w", err)
	}

	return comments, nil
}
