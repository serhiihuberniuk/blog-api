package service

import (
	"context"

	"github.com/serhiihuberniuk/blog-api/models"
)

type Service struct {
	repo repository
}

type repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error

	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, postID string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, post *models.Post) error
	ListPosts(ctx context.Context, pagination models.Pagination,
		filter models.FilterPosts, sort models.SortPosts) (*[]models.Post, error)

	CreateComment(ctx context.Context, comment *models.Comment) error
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, comment *models.Comment) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) (*[]models.Comment, error)
}
