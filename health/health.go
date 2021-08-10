package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type healthStatus struct {
	PostgresDb string `json:"postgres_db"`
	MongoDb    string `json:"mongo_db"`
}

type HandlerHealth struct {
	PostgresDb *pgxpool.Pool
	MongoDB    bool
}

func writeJsonResponse(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *HandlerHealth) HealthRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", h.Health)

	return router
}

func (h *HandlerHealth) Health(w http.ResponseWriter, r *http.Request) {
	status := &healthStatus{}
	codes := []int{
		h.postgresHealth(r, status),
		h.mongoHealth(r, status),
	}

	data, _ := json.Marshal(status)

	for _, code := range codes {
		if code != http.StatusOK {
			writeJsonResponse(w, http.StatusInternalServerError, data)
			return
		}
	}

	writeJsonResponse(w, http.StatusOK, data)
}

func (h *HandlerHealth) postgresHealth(r *http.Request, status *healthStatus) int {
	if err := h.PostgresDb.Ping(r.Context()); err != nil {
		status.PostgresDb = fmt.Sprintf("error: %v", err)
		return http.StatusInternalServerError
	}

	status.PostgresDb = "OK"

	return http.StatusOK
}

func (h *HandlerHealth) mongoHealth(r *http.Request, status *healthStatus) int {
	if !h.MongoDB {
		status.MongoDb = fmt.Sprintf("error:")
		return http.StatusInternalServerError
	}

	status.MongoDb = "OK"

	return http.StatusOK
}
