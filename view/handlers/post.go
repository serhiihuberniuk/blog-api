package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/models"
)

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "cannot decode data from JSON", http.StatusBadRequest)

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

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

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

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		http.Error(w, "cannot find user with such ID", http.StatusNotFound)

		return
	}

	in := viewmodels.UpdatePostRequest{
		Title:       post.Title,
		Description: post.Description,
		Tags:        post.Tags,
	}

	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "cannot decode data from JSON", http.StatusBadRequest)
	}

	err = h.service.UpdatePost(r.Context(), models.UpdatePostPayload{
		PostID:      postID,
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
	})
	if err != nil {
		http.Error(w, "cannot update post", http.StatusBadRequest)

		return
	}

	post, err = h.service.GetPost(r.Context(), postID)
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

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

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
	vars := r.URL.Query()
	filterByField := vars.Get("filterByField")

	filterValue := vars.Get("filterValue")
	if filterValue == "" {
		filterByField = ""
	}

	sortByField := vars.Get("sortByField")
	isAscString := vars.Get("isAsc")

	isAsc := true
	if isAscString == "false" {
		isAsc = false
	}

	limitString := vars.Get("limit")

	limit, err := strconv.Atoi(limitString)
	if err != nil && limitString != "" {
		http.Error(w, "cannot convert limit into uint", http.StatusBadRequest)

		return
	}

	if limit == 0 || limit > maxlimit {
		limit = maxlimit
	}

	offsetString := vars.Get("offset")

	offset, err := strconv.Atoi(offsetString)
	if err != nil && offsetString != "" {
		http.Error(w, "cannot convert offset into uint", http.StatusBadRequest)

		return
	}

	posts, err := h.service.ListPosts(r.Context(), models.Pagination{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	},
		models.FilterPosts{
			Field: models.FilterPostsByField(filterByField),
			Value: filterValue,
		},
		models.SortPosts{
			SortByField: models.SortPostsByField(sortByField),
			IsASC:       isAsc,
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

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(outs); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}
