package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreateComment(ctx context.Context, comment *models.Comment) error {
	const sql = "INSERT INTO comments (id, content, created_by, created_at, post_id) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.Db.Exec(ctx, sql, comment.ID, comment.Content, comment.CreatedBy, comment.CreatedAt, comment.PostID)
	if err != nil {
		return fmt.Errorf("cannot create comment, %w", err)
	}

	return nil
}

func (r *Repository) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	const sql = "SELECT id, content, created_by, created_at, post_id FROM comments WHERE id=$1"

	var comment models.Comment

	err := pgxscan.Get(ctx, r.Db, &comment, sql, commentID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("cannot get comment, %w", err)
	}

	return &comment, nil
}

func (r *Repository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	const sql = "UPDATE comments SET content=$1 WHERE id=$2"

	result, err := r.Db.Exec(ctx, sql, comment.Content, comment.ID)
	if err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID string) error {
	const sql = "DELETE FROM comments WHERE id=$1"

	result, err := r.Db.Exec(ctx, sql, commentID)
	if err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) ListComments(ctx context.Context, pagination models.Pagination,
	filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error) {
	query := squirrel.Select("*").
		From("comments")
	if filter.Field != "" {
		query = query.Where(fmt.Sprintf("%s=$1", filter.Field), filter.Value)
	}

	if sort.Field != "" {
		if !sort.IsASC {
			query = query.OrderBy(string(sort.Field) + " DESC")
		} else {
			query = query.OrderBy(string(sort.Field))
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
		return nil, fmt.Errorf("cannot get list of commentss: %w", err)
	}

	var comments []*models.Comment

	err = pgxscan.Select(ctx, r.Db, &comments, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of commentss: %w", err)
	}

	return comments, nil
}
