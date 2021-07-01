package service

import (
	"context"
	"github.com/serhiihuberniuk/blog-api/models"
	"time"
)

var now = time.Now()

type Service struct {
	repo Repository
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userId string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error

	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, postId string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, post *models.Post) error
	//ListPosts(ctx context.Context, posts ...models.Post) ([]models.Post, error)

	CreateComment(ctx context.Context, comment *models.Comment) error
	GetComment(ctx context.Context, commentId string) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, comment *models.Comment) error
	//ListComments(ctx context.Context, comments ...models.Comment) ([]models.Comment, error)
}
