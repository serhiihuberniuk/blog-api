package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
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

func handleError(err error) int {
	if errors.Is(err, models.ErrNotFoundUser) {
		return http.StatusNotFound
	}

	if errors.Is(err, models.ErrNotFoundPost) {
		return http.StatusNotFound
	}

	if errors.Is(err, models.ErrNotFoundComment) {
		return http.StatusNotFound
	}

	if errors.As(err, &validation.Errors{}) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func decodeFromJson(w http.ResponseWriter, r *http.Request, a interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		http.Error(w, "cannot decode data from JSON", http.StatusBadRequest)

		return false
	}

	return true
}

func encodeIntoJson(w http.ResponseWriter, a interface{}) bool {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(a); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return false
	}

	return true
}

func GetQueryParams(r *http.Request, allowFilterFn, allowSortFn func(string) bool) (*queryParam, error) {
	vars := r.URL.Query()

	var err error

	queryParams := queryParam{}

	if v := vars.Get("filter-field"); v != "" && allowFilterFn(v) {
		queryParams.filterByField = v
		queryParams.filterValue = vars.Get("filter-value")
	}

	if v := vars.Get("sort-field"); v != "" && allowSortFn(v) {
		queryParams.sortByField = v
	}

	if vars.Get("is-asc") != "false" {
		queryParams.isAsc = true
	}

	queryParams.limit, err = strconv.Atoi(vars.Get("limit"))
	if err != nil && vars.Get("limit") != "" {
		return nil, fmt.Errorf("cannot convert limit into int: %w", err)
	}

	if queryParams.limit <= 0 || queryParams.limit > maxLimit {
		queryParams.limit = maxLimit
	}

	queryParams.offset = 0
	if vars.Get("offset") != "" {
		queryParams.offset, err = strconv.Atoi(vars.Get("offset"))
		if err != nil && vars.Get("offset") != "" {
			return nil, fmt.Errorf("cannot convert offset into int: %w", err)
		}
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

func NewRestHandlers(s service) *Handlers {
	return &Handlers{
		service: s,
	}
}
