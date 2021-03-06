package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/rest/models"
)

var (
	allowedFilterCommentsFields = map[string]models.FilterCommentsByField{
		string(models.FilterCommentsByPost):      models.FilterCommentsByPost,
		string(models.FilterCommentsByCreatedAt): models.FilterCommentsByCreatedAt,
		string(models.FilterCommentsByAuthor):    models.FilterCommentsByAuthor,
	}

	allowedSortCommentsFields = map[string]models.SortCommentsByField{
		string(models.SortCommentByCreatedAt): models.SortCommentByCreatedAt,
	}
)

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreateCommentRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	commentID, err := h.service.CreateComment(r.Context(), models.CreateCommentPayload{
		Content: in.Content,
		PostID:  in.PostID,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) GetComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	var in viewmodels.UpdateCommentRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	err := h.service.UpdateComment(r.Context(), models.UpdateCommentPayload{
		CommentID: commentID,
		Content:   in.Content,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	comment, err := h.service.GetComment(r.Context(), commentID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	if err := h.service.DeleteComment(r.Context(), commentID); err != nil {
		errorStatusHttp(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetListOfComments(w http.ResponseWriter, r *http.Request) {
	queryParams, err := GetQueryParams(r, func(s string) bool {
		_, ok := allowedFilterCommentsFields[s]

		return ok
	}, func(s string) bool {
		_, ok := allowedSortCommentsFields[s]

		return ok
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	comments, err := h.service.ListComments(r.Context(), models.Pagination{
		Limit:  uint64(queryParams.limit),
		Offset: uint64(queryParams.offset),
	}, models.FilterComments{
		Field: models.FilterCommentsByField(queryParams.filterByField),
		Value: queryParams.filterValue,
	}, models.SortComments{
		Field: models.SortCommentsByField(queryParams.sortByField),
		IsASC: queryParams.isAsc,
	})
	if err != nil {
		errorStatusHttp(w, err)

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

	if !encodeIntoJson(w, outs) {
		return
	}
}
