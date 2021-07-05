package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (r *Repository) CreateComment(ctx context.Context, comment *models.Comment) error {
	sql, args, err := squirrel.Insert("comments").
		Values(comment.ID, comment.Content, comment.CreatedBy, comment.CreatedAt, comment.PostID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create comment, %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot create comment, %w", err)
	}

	return nil
}

func (r *Repository) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	var comment models.Comment

	sql, args, err := squirrel.Select("*").
		From("comments").
		Where("id=$1", commentID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cennot get comment, %w", err)
	}

	err = r.Db.QueryRow(ctx, sql, args...).
		Scan(
			&comment.ID,
			&comment.Content,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&comment.PostID,
		)
	if err != nil {
		return nil, fmt.Errorf("cennot get comment, %w", err)
	}

	return &comment, nil
}

func (r *Repository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	sql, args, err := squirrel.Update("comments").
		Set("content", comment.Content).
		Where("id=$1", comment.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, comment *models.Comment) error {
	sql, args, err := squirrel.Delete("comments").
		Where("id=$1", comment.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}

	return nil
}

func (r *Repository) ListComments(ctx context.Context, pagination models.Pagination,
	filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error) {
	var comments []*models.Comment

	query := squirrel.Select("*").
		From("comments")
	if pred := fmt.Sprintf("%s=$1", filter.Field); filter.Field != "" {
		query = query.Where(pred, filter.Value)
	}

	if sort.Field != "" {
		if !sort.IsASC {
			query = query.OrderBy(string(sort.Field) + " DESC")
		}

		query = query.OrderBy(string(sort.Field))
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

	err = pgxscan.Select(ctx, r.Db, &comments, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of commentss: %w", err)
	}

	return comments, nil
}
