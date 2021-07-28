package graph

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
const (
	maxLimit = 50
)

var (
	allowedFilterPostsFields = map[model.FilterPostsField]models.FilterPostsByField{
		model.FilterPostsFieldCreatedBy: models.FilterPostsByCreatedBy,
		model.FilterPostsFieldTags:      models.FilterPostsByTags,
		model.FilterPostsFieldTitle:     models.FilterPostsByTitle,
	}
	allowedSortPostsFields = map[model.SortPostsField]models.SortPostsByField{
		model.SortPostsFieldCreatedAt: models.SortPostsByCreatedAt,
		model.SortPostsFieldTitle:     models.SortPostsByTitle,
	}
	allowedFilterCommentsFields = map[model.FilterCommentsField]models.FilterCommentsByField{
		model.FilterCommentsFieldCreatedAt: models.FilterCommentsByCreatedAt,
		model.FilterCommentsFieldCreatedBy: models.FilterCommentsByAuthor,
		model.FilterCommentsFieldPostID:    models.FilterCommentsByPost,
	}
	allowedSortCommentsFields = map[model.SortCommentsField]models.SortCommentsByField{
		model.SortCommentsFieldCreatedAt: models.SortCommentByCreatedAt,
	}
)

type Resolver struct {
	service
}

func getPaginationParams(paginationInput *model.PaginationInput) models.Pagination {
	pagination := models.Pagination{
		Limit:  maxLimit,
		Offset: 0,
	}

	if paginationInput != nil {
		if paginationInput.Limit != nil {
			pagination.Limit = uint64(*paginationInput.Limit)
		}

		if pagination.Limit > maxLimit {
			pagination.Limit = maxLimit
		}

		if paginationInput.Offset != nil {
			pagination.Offset = uint64(*paginationInput.Offset)
		}
	}

	return pagination
}

func getUser(ctx context.Context, r *Resolver, id string) (*model.User, error) {
	user, err := r.service.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func getPost(ctx context.Context, r *Resolver, id string) (*model.Post, error) {
	post, err := r.service.GetPost(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	postAuthor, err := getUser(ctx, r, post.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of post: %w", err)
	}

	return &model.Post{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   postAuthor,
		CreatedAt:   post.CreatedAt.String(),
		Tags:        post.Tags,
	}, nil
}

func getComment(ctx context.Context, r *Resolver, id string) (*model.Comment, error) {
	comment, err := r.service.GetComment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get comment, %w", err)
	}

	author, err := getUser(ctx, r, comment.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of comment: %w", err)
	}

	post, err := getPost(ctx, r, comment.PostID)
	if err != nil {
		return nil, fmt.Errorf("cannot get parent post: %w", err)
	}

	return &model.Comment{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedBy: author,
		CreatedAt: comment.CreatedAt.String(),
		Post:      post,
	}, nil
}

type service interface {
	CreateUser(ctx context.Context, payload models.CreateUserPayload) (string, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, payload models.UpdateUserPayload) error
	DeleteUser(ctx context.Context, userID string) error

	CreatePost(ctx context.Context, payload models.CreatePostPayload) (string, error)
	GetPost(ctx context.Context, postID string) (*models.Post, error)
	UpdatePost(ctx context.Context, payload models.UpdatePostPayload) error
	DeletePost(ctx context.Context, postID string) error
	ListPosts(ctx context.Context, pagination models.Pagination,
		filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error)

	CreateComment(ctx context.Context, payload models.CreateCommentPayload) (string, error)
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, payload models.UpdateCommentPayload) error
	DeleteComment(ctx context.Context, commentId string) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error)
}

func NewResolver(s service) *Resolver {
	return &Resolver{
		s,
	}
}
