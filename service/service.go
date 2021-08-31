package service

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/serhiihuberniuk/blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Service struct {
	repo                           repository
	privateKey                     *rsa.PrivateKey
	currentUserInformationProvider currentUserInformationProvider
}

type currentUserInformationProvider interface {
	GetCurrentUserID(ctx context.Context) string
}

type repository interface {
	Login(ctx context.Context, email string) (*models.User, error)

	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, postID string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, postID string) error
	ListPosts(ctx context.Context, pagination models.Pagination,
		filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error)

	CreateComment(ctx context.Context, comment *models.Comment) error
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, commentID string) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error)
}

func NewService(r repository, privateKey []byte, p currentUserInformationProvider) (*Service, error) {
	privateRSA, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error occurred while parsing private key: %w", err)
	}

	return &Service{
		repo:                           r,
		privateKey:                     privateRSA,
		currentUserInformationProvider: p,
	}, nil
}

func generateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error occurred while hashing the password, %w", err)
	}

	return string(hashedPassword), nil
}
