package health

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerHealth struct {
	checks []check
}

type check func(ctx context.Context) error

func NewHandlerHealth(c ...check) *HandlerHealth {
	return &HandlerHealth{
		c,
	}
}

func (h *HandlerHealth) HealthRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", h.Health)

	return router
}

func (h *HandlerHealth) Health(w http.ResponseWriter, r *http.Request) {
	for _, check := range h.checks {
		if err := check(r.Context()); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
