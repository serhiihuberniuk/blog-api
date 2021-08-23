package grpcHandlers

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func errorStatusGrpc(err error) error {
	if errors.Is(err, models.ErrNotFound) {
		return status.Error(codes.NotFound, codes.NotFound.String())
	}

	if errors.Is(err, models.ErrNotAuthenticated) {
		return status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	if errors.As(err, &validation.Errors{}) {
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	return status.Error(codes.Internal, codes.Internal.String())
}

type Handlers struct {
	service                        service
	currentUserInformationProvider currentUserInformationProvider
}

type currentUserInformationProvider interface {
	GetCurrentUserID(ctx context.Context) string
}

type service interface {
	Login(ctx context.Context, payload models.LoginPayload) (string, error)

	CreateUser(ctx context.Context, payload models.CreateUserPayload) (string, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, payload models.UpdateUserPayload) error
	DeleteUser(ctx context.Context) error

	CreatePost(ctx context.Context, payload models.CreatePostPayload) (string, error)
	GetPost(ctx context.Context, postID string) (*models.Post, error)
	UpdatePost(ctx context.Context, payload models.UpdatePostPayload) error
	DeletePost(ctx context.Context, postID string) error
	ListPosts(ctx context.Context, pagination models.Pagination,
		filter models.FilterPosts, sort models.SortPosts) ([]*models.Post, error)

	CreateComment(ctx context.Context, payload models.CreateCommentPayload) (string, error)
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, payload models.UpdateCommentPayload) error
	DeleteComment(ctx context.Context, commentID string) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error)
}

func NewGrpcHandlers(s service, p currentUserInformationProvider) *Handlers {
	return &Handlers{
		service:                        s,
		currentUserInformationProvider: p,
	}
}
