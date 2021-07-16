package handlers

import (
	models2 "github.com/serhiihuberniuk/blog-api/view/rest/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in models2.CreateUserRequest

	if !decodeFromJson(w, r, in) {
		return
	}

	userID, err := h.service.CreateUser(r.Context(), models.CreateUserPayload{
		Name:  in.Name,
		Email: in.Email,
	})
	if err != nil {
		http.Error(w, "cannot create user", http.StatusBadRequest)

		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "cannot get created user", http.StatusInternalServerError)

		return
	}

	out := models2.GetUserResponse{
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
		http.Error(w, "cannot find user with such ID", http.StatusNotFound)

		return
	}

	out := models2.GetUserResponse{
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

	var in models2.UpdateUserRequest

	if !decodeFromJson(w, r, in) {
		return
	}

	err := h.service.UpdateUser(r.Context(), models.UpdateUserPayload{
		UserID: userID,
		Name:   in.Name,
		Email:  in.Email,
	})
	if err != nil {
		http.Error(w, "cannot update user", http.StatusBadRequest)

		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "cannot get updated user", http.StatusNotFound)

		return
	}

	out := models2.GetUserResponse{
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
		http.Error(w, "cannot find user with such ID to delete", http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)
}
