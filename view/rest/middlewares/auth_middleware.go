package middlewares

import (
	"context"
	"net/http"
	"strings"
)

const (
	authorizationHeader  = "Authorization"
	bearerAuthentication = "bearer"
)

type AuthMiddleware struct {
	service                        service
	currentUserInformationProvider currentUserInformationProvider
}

func NewAuthMiddleware(s service, p currentUserInformationProvider) *AuthMiddleware {
	return &AuthMiddleware{
		service:                        s,
		currentUserInformationProvider: p,
	}
}

type currentUserInformationProvider interface {
	SetCurrentUserID(ctx context.Context, userID string) context.Context
	GetCurrentUserID(ctx context.Context) string
}

type service interface {
	ParseToken(tokenString string) (string, error)
}

func (m *AuthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
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

		userID, err := m.service.ParseToken(headerSplit[1])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		ctx := m.currentUserInformationProvider.SetCurrentUserID(r.Context(), userID)

		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) GetCurrentUserID(ctx context.Context) string {
	return m.currentUserInformationProvider.GetCurrentUserID(ctx)
}
