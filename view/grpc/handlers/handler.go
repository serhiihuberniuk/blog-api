package grpcHandlers

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"google.golang.org/grpc/codes"
)

const maxLimit = 50

func getPaginationParam(p *pb.Pagination) models.Pagination {
	pagination := models.Pagination{}

	if p != nil {
		limit := p.GetLimit()
		if limit <= 0 || limit > maxLimit {
			limit = maxLimit
		}

		offset := p.GetOffset()
		if offset < 0 {
			offset = 0
		}

		pagination = models.Pagination{
			Limit:  uint64(limit),
			Offset: uint64(offset),
		}
	}

	return pagination
}

func handleError(err error) (codes.Code, error) {
	if errors.Is(err, models.ErrNotFoundUser) {
		return codes.NotFound, models.ErrNotFoundUser
	}

	if errors.Is(err, models.ErrNotFoundPost) {
		return codes.NotFound, models.ErrNotFoundPost
	}

	if errors.Is(err, models.ErrNotFoundComment) {
		return codes.NotFound, models.ErrNotFoundComment
	}

	if errors.As(err, &validation.Errors{}) {
		return codes.InvalidArgument, err
	}

	return codes.Internal, err
}

type Handlers struct {
	service
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

func NewGrpcHandlers(s service) *Handlers {
	return &Handlers{
		service: s,
	}
}
