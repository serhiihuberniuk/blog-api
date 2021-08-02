package handlers

import (
	"errors"
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
		if errors.Is(err, models.ErrorBadRequest) {
			http.Error(w, models.ErrorBadRequest.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "cannot create user", http.StatusInternalServerError)

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
		if errors.Is(err, models.UserNotFound) {
			http.Error(w, models.UserNotFound.Error(), http.StatusNotFound)
		}

		http.Error(w, "cannot get user", http.StatusInternalServerError)

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
		if errors.Is(err, models.UserNotFound) {
			http.Error(w, models.UserNotFound.Error(), http.StatusNotFound)

			return
		}

		if errors.Is(err, models.ErrorBadRequest) {
			http.Error(w, models.ErrorBadRequest.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "cannot update user", http.StatusBadRequest)

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
		if errors.Is(err, models.UserNotFound) {
			http.Error(w, models.UserNotFound.Error(), http.StatusNotFound)

			return
		}

		http.Error(w, "cannot delete user", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
