package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/serhiihuberniuk/blog-api/models"
)

const maxLimit = 50

type queryParam struct {
	filterByField string
	filterValue   string
	sortByField   string
	isAsc         bool
	limit         int
	offset        int
}

func decodeFromJson(w http.ResponseWriter, r *http.Request, a interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "cannot decode data from JSON", http.StatusBadRequest)

		return fmt.Errorf("error occurred while decoding from JSON: %w", err)
	}

	return nil
}

func encodeIntoJson(w http.ResponseWriter, a interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(a); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return fmt.Errorf("error occurred while encoding into JSON: %w", err)
	}

	return nil
}

func SetQueryParams(r *http.Request) (*queryParam, error) {
	vars := r.URL.Query()

	var err error

	queryParams := queryParam{
		filterByField: vars.Get("filter-field"),
		filterValue:   vars.Get("filter-value"),
		sortByField:   vars.Get("sort-field"),
		isAsc:         true,
		limit:         maxLimit,
		offset:        0,
	}

	if queryParams.filterValue == "" {
		queryParams.filterByField = ""
	}

	if vars.Get("is-asc") == "false" {
		queryParams.isAsc = false
	}

	queryParams.limit, err = strconv.Atoi(vars.Get("limit"))
	if err != nil && vars.Get("limit") != "" {
		return nil, fmt.Errorf("cannot convert limit into int: %w", err)
	}

	if queryParams.limit <= 0 || queryParams.limit > maxLimit {
		queryParams.limit = maxLimit
	}

	offsetString := vars.Get("offset")

	queryParams.offset, err = strconv.Atoi(offsetString)
	if err != nil && offsetString != "" {
		return nil, fmt.Errorf("cannot convert offset into int: %w", err)
	}

	if queryParams.offset < 0 {
		queryParams.offset = 0
	}

	return &queryParams, nil
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

type Handlers struct {
	service
}

func NewHandlers(s service) *Handlers {
	return &Handlers{
		service: s,
	}
}
