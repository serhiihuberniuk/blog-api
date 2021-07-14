package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/models"
)

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreateCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "cannot encode data from JSON", http.StatusBadRequest)

		return
	}

	commentID, err := h.service.CreateComment(r.Context(), models.CreateCommentPayload{
		Content:  in.Content,
		PostID:   in.PostID,
		AuthorID: in.AuthorID,
	})
	if err != nil {
		http.Error(w, "cannot create comment", http.StatusBadRequest)

		return
	}

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		http.Error(w, "cannot get created comment", http.StatusNotFound)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) GetComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		http.Error(w, "cannot get comment with such ID", http.StatusNotFound)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		http.Error(w, "cannot get comment with such ID", http.StatusNotFound)

		return
	}

	in := viewmodels.UpdateCommentRequest{
		Content: comment.Content,
	}

	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "cannot decode data from JSON", http.StatusBadRequest)

		return
	}

	err = h.service.UpdateComment(r.Context(), models.UpdateCommentPayload{
		CommentID: commentID,
		Content:   in.Content,
	})
	if err != nil {
		http.Error(w, "cannot update comment", http.StatusBadRequest)

		return
	}

	comment, err = h.service.GetComment(r.Context(), commentID)
	if err != nil {
		http.Error(w, "cannot get updated comment", http.StatusInternalServerError)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	if err := h.service.DeleteComment(r.Context(), commentID); err != nil {
		http.Error(w, "cannot delete comment with such ID", http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetListOfComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := h.service.ListComments(r.Context(), models.Pagination{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}, models.FilterComments{
		Field: models.FilterCommentsByField(filterByField),
		Value: filterValue,
	}, models.SortComments{
		Field: models.SortCommentsByField(sortByField),
		IsASC: isAsc,
	})
	if err != nil {
		http.Error(w, "cannot get comments", http.StatusBadRequest)

		return
	}

	outs := make([]viewmodels.GetCommentResponse, 0)

	for _, comment := range comments {
		out := viewmodels.GetCommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			CreatedBy: comment.CreatedBy,
			PostID:    comment.PostID,
		}

		outs = append(outs, out)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(outs); err != nil {
		http.Error(w, "cannot encode data into JSON", http.StatusInternalServerError)

		return
	}
}
