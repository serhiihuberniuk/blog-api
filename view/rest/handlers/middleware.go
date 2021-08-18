package handlers

import (
	"net/http"
	"strings"
)

const (
	authorizationHeader  = "Authorization"
	bearerAuthentication = "bearer"
)

func (h *Handlers) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		headerSplit := strings.Split(header, " ")
		if len(headerSplit) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		if strings.ToLower(headerSplit[0]) != bearerAuthentication {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		userID, err := h.service.ParseToken(headerSplit[1])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		ctx := h.provider.SetCurrentUserID(r.Context(), userID)

		next(w, r.WithContext(ctx))
	}
}
