package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreatePost(ctx context.Context, post *models.Post) error {
	const sql = "INSERT INTO posts (id, title, description, created_by, created_at, tags) VALUES ($1,$2,$3,$4,$5,$6)"

	_, err := r.Db.Exec(ctx, sql, post.ID, post.Title, post.Description, post.CreatedBy, post.CreatedAt, post.Tags)
	if err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	return nil
}

func (r *Repository) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	const sql = "SELECT id, title, description, created_by, created_at, tags FROM posts WHERE id=$1"

	var post models.Post

	err := pgxscan.Get(ctx, r.Db, &post, sql, postID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, models.ErrNotFoundPost
		}

		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return &post, nil
}

func (r *Repository) UpdatePost(ctx context.Context, post *models.Post) error {
	const sql = "UPDATE posts SET title=$1, description=$2, tags=$3 WHERE id=$4"

	result, err := r.Db.Exec(ctx, sql, post.Title, post.Description, post.Tags, post.ID)
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFoundPost
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, postID string) error {
	const sql = "DELETE FROM posts WHERE id=$1"

	result, err := r.Db.Exec(ctx, sql, postID)
	if err != nil {
		return fmt.Errorf("cannot delete post, %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFoundPost
	}

	return nil
}

func (r *Repository) ListPosts(ctx context.Context, pagination models.Pagination,
	filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error) {
	query := squirrel.Select("*").
		From("posts")

	if filter.Field != "" {
		if filter.Field == models.FilterPostsByTags {
			query = query.Where(fmt.Sprintf("%s::jsonb ? $1", filter.Field), filter.Value)
		} else {
			query = query.Where(fmt.Sprintf("%s=$1", filter.Field), filter.Value)
		}
	}

	if sort.SortByField != "" {
		if !sort.IsASC {
			query = query.OrderBy(string(sort.SortByField) + " DESC")
		} else {
			query = query.OrderBy(string(sort.SortByField))
		}
	}

	if pagination.Offset != 0 {
		query = query.Offset(pagination.Offset)
	}

	if pagination.Limit != 0 {
		query = query.Limit(pagination.Limit)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	var posts []*models.Post

	err = pgxscan.Select(ctx, r.Db, &posts, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	return posts, nil
}
