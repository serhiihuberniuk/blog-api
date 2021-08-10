package health

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerHealth struct {
	repo
}

type repo interface {
	HealthCheck(ctx context.Context) error
}

func NewHandlerHealth(r repo) *HandlerHealth {
	return &HandlerHealth{
		r,
	}
}

func (h *HandlerHealth) HealthRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", h.Health)

	return router
}

func (h *HandlerHealth) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.HealthCheck(r.Context()); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
