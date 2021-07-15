package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/models"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.CreateUserRequest

	if err := decodeFromJson(w, r, in); err != nil {
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

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err = encodeIntoJson(w, out); err != nil {
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

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err = encodeIntoJson(w, out); err != nil {
		return
	}
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	var in viewmodels.UpdateUserRequest

	if err := decodeFromJson(w, r, in); err != nil {
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

	out := viewmodels.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err = encodeIntoJson(w, out); err != nil {
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
