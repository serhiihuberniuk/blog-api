package handlers

import (
	"net/http"

	"github.com/serhiihuberniuk/blog-api/models"
	viewmodels "github.com/serhiihuberniuk/blog-api/view/rest/models"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var in viewmodels.LoginRequest

	if !decodeFromJson(w, r, &in) {
		return
	}

	token, err := h.service.Login(r.Context(), models.LoginPayload{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		errorStatusHttp(w, err)
	}

	out := &viewmodels.LoginResponse{
		Token: token,
	}

	if !encodeIntoJson(w, out) {
		return
	}
}
