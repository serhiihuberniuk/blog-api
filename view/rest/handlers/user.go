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
		Name:  in.Name,
		Email: in.Email,
	})
	if err != nil {
		code := handleError(err)
		http.Error(w, "cannot create user", code)

		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "cannot get created user", http.StatusInternalServerError)

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
		code := handleError(err)
		http.Error(w, "cannot get user", code)

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
	userID := mux.Vars(r)["id"]

	var in viewmodels.UpdateUserRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	err := h.service.UpdateUser(r.Context(), models.UpdateUserPayload{
		UserID: userID,
		Name:   in.Name,
		Email:  in.Email,
	})
	if err != nil {
		code := handleError(err)
		http.Error(w, "cannot update user", code)

		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "cannot get updated user", http.StatusInternalServerError)

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
	userID := mux.Vars(r)["id"]

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		code := handleError(err)
		http.Error(w, "cannot delete user", code)

		return
	}

	w.WriteHeader(http.StatusOK)
}
