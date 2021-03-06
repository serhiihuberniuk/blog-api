package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/generated"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/model"
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
	service                        service
	currentUserInformationProvider currentUserInformationProvider
}

type currentUserInformationProvider interface {
	GetCurrentUserID(ctx context.Context) string
}

func getPaginationParams(paginationInput *model.PaginationInput) models.Pagination {
	pagination := models.Pagination{
		Limit:  maxLimit,
		Offset: 0,
	}

	if paginationInput != nil {
		pagination = models.Pagination{
			Limit:  uint64(paginationInput.Limit),
			Offset: uint64(paginationInput.Offset),
		}

		if pagination.Limit > maxLimit {
			pagination.Limit = maxLimit
		}
	}

	return pagination
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
	DeleteComment(ctx context.Context, commentId string) error
	ListComments(ctx context.Context, pagination models.Pagination,
		filter models.FilterComments, sort models.SortComments) ([]*models.Comment, error)
}

func NewResolverConfig(s service, p currentUserInformationProvider) generated.Config {
	r := &Resolver{
		service:                        s,
		currentUserInformationProvider: p,
	}

	resolverConfig := generated.Config{
		Resolvers: r,
	}

	resolverConfig.Directives.IsAuthenticated = func(ctx context.Context, obj interface{},
		next graphql.Resolver) (res interface{}, err error) {
		ctxUserID := r.currentUserInformationProvider.GetCurrentUserID(ctx)
		if ctxUserID != "" {
			return next(ctx)
		}

		return nil, models.ErrNotAuthenticated
	}

	return resolverConfig
}
