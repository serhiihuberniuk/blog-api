package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/rest/models"
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
		AuthorID:    in.AuthorID,
	})
	if err != nil {
		http.Error(w, "cannot create post", http.StatusInternalServerError)

		return
	}

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		http.Error(w, "cannot get created post", http.StatusNotFound)

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
		http.Error(w, "cannot find post with such ID", http.StatusNotFound)

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
		http.Error(w, "cannot update post", http.StatusBadRequest)

		return
	}

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		http.Error(w, "cannot get updated post", http.StatusInternalServerError)

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
		http.Error(w, "cannot get post with such ID", http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetListOfPosts(w http.ResponseWriter, r *http.Request) {
	queryParams, err := GetQueryParams(r)
	if err != nil {
		http.Error(w, "requested parameters is bad", http.StatusBadRequest)

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
		http.Error(w, "cannot get posts", http.StatusBadRequest)

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
