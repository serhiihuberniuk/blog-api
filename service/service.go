package service

import (
	"context"

	"github.com/serhiihuberniuk/blog-api/models"
)

type Service struct {
	repo Repository
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User, payload CreateUserPayload) error
	GetUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, user *models.User, payload UpdateUserPayload) error
	DeleteUser(ctx context.Context, user models.User) error

	CreatePost(ctx context.Context, post *models.Post, payload CreatePostPayload, user models.User) error
	GetPost(ctx context.Context, post models.Post) (models.Post, error)
	UpdatePost(ctx context.Context, post models.Post, payload UpdatePostPayload) error
	DeletePost(ctx context.Context, post models.Post) error
	ListPosts(ctx context.Context, posts ...models.Post) ([]models.Post, error)

	CreateComment(ctx context.Context, payload CreateCommentPayload, comment *models.Comment, user models.User, post models.Post) error
	GetComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment, payload UpdateCommentPayload)
	DeleteComment(ctx context.Context, comment models.Comment) error
	ListComments(ctx context.Context, comments ...models.Comment) ([]models.Comment, error)
}
