package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/rest/models"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreateUserRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	userID, err := h.service.CreateUser(r.Context(), models.CreateUserPayload{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.UpdateUserRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	err := h.service.UpdateUser(r.Context(), models.UpdateUserPayload{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	user, err := h.service.GetUser(r.Context(), h.currentUserInformationProvider.GetCurrentUserID(r.Context()))
	if err != nil {
		errorStatusHttp(w, err)

		return
	}

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteUser(r.Context()); err != nil {
		errorStatusHttp(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}
