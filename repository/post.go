package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreatePost(ctx context.Context, post *models.Post) error {
	sql, args, err := squirrel.Insert("posts").
		Values(post.ID, post.Title, post.Description, post.CreatedAt, post.CreatedBy, post.Tags).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot create post: %w", err)
	}

	return nil
}

func (r *Repository) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	var post models.Post

	sql, args, err := squirrel.Select("*").
		From("posts").
		Where("id=$1", postID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot get post, %w", err)
	}

	err = r.Db.QueryRow(ctx, sql, args...).
		Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.CreatedAt,
			&post.CreatedBy,
			&post.Tags,
		)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return &post, nil
}

func (r *Repository) UpdatePost(ctx context.Context, post *models.Post) error {
	sql, args, err := squirrel.Update("posts").
		Set("title", post.Title).
		Set("description", post.Description).
		Set("tags", post.Tags).
		Where("id=$1", post.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, post *models.Post) error {
	sql, args, err := squirrel.Delete("posts").
		Where("id=$1", post.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot delete post, %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete post, %w", err)
	}

	return nil
}

func (r *Repository) ListPosts(ctx context.Context, pagination models.Pagination,
	filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error) {
	var posts []*models.Post

	query := squirrel.Select("*").
		From("posts")
	if pred := fmt.Sprintf("%s=$1", filter.Field); filter.Field != "" {
		query = query.Where(pred, filter.Value)
	}

	if sort.SortByField != "" {
		if !sort.IsASC {
			query = query.OrderBy(string(sort.SortByField) + " DESC")
		}

		query = query.OrderBy(string(sort.SortByField))
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

	err = pgxscan.Select(ctx, r.Db, &posts, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	return posts, nil
}
