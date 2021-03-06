package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/rest/models"
)

var (
	allowedFilterPostsFields = map[string]models.FilterPostsByField{
		string(models.FilterPostsByCreatedBy): models.FilterPostsByCreatedBy,
		string(models.FilterPostsByTitle):     models.FilterPostsByTitle,
		string(models.FilterPostsByTags):      models.FilterPostsByTags,
	}
	allowedSortPostsFields = map[string]models.SortPostsByField{
		string(models.SortPostsByCreatedAt): models.SortPostsByCreatedAt,
		string(models.SortPostsByTitle):     models.SortPostsByTitle,
	}
)

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreatePostRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	postID, err := h.service.CreatePost(r.Context(), models.CreatePostPayload{
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		CreatedBy:   post.CreatedBy,
		Tags:        post.Tags,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		CreatedBy:   post.CreatedBy,
		Tags:        post.Tags,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	var in viewmodels.UpdatePostRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	err := h.service.UpdatePost(r.Context(), models.UpdatePostPayload{
		PostID:      postID,
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		CreatedBy:   post.CreatedBy,
		Tags:        post.Tags,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	if err := h.service.DeletePost(r.Context(), postID); err != nil {
		errorStatusHttp(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetListOfPosts(w http.ResponseWriter, r *http.Request) {
	queryParams, err := GetQueryParams(r, func(s string) bool {
		_, ok := allowedFilterPostsFields[s]

		return ok
	}, func(s string) bool {
		_, ok := allowedSortPostsFields[s]

		return ok
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	posts, err := h.service.ListPosts(r.Context(), models.Pagination{
		Limit:  uint64(queryParams.limit),
		Offset: uint64(queryParams.offset),
	},
		models.FilterPosts{
			Field: models.FilterPostsByField(queryParams.filterByField),
			Value: queryParams.filterValue,
		},
		models.SortPosts{
			SortByField: models.SortPostsByField(queryParams.sortByField),
			IsASC:       queryParams.isAsc,
		})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	outs := make([]viewmodels.GetPostResponse, 0)

	for _, post := range posts {
		out := viewmodels.GetPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			CreatedAt:   post.CreatedAt,
			CreatedBy:   post.CreatedBy,
			Tags:        post.Tags,
		}

		outs = append(outs, out)
	}

	if !encodeIntoJson(w, outs) {
		return
	}
}
