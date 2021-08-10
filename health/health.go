package health

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerHealth struct {
	service []service
}

type service interface {
	HealthCheck(ctx context.Context) error
}

func NewHandlerHealth(s ...service) *HandlerHealth {
	return &HandlerHealth{
		s,
	}
}

func (h *HandlerHealth) HealthRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", h.Health)

	return router
}

func (h *HandlerHealth) Health(w http.ResponseWriter, r *http.Request) {
	for _, service := range h.service {
		if err := service.HealthCheck(r.Context()); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
