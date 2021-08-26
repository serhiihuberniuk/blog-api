package graphqlMiddlewares

import (
	"context"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
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

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)

		if header == "" {
			next.ServeHTTP(w, r)

			return
		}

		userID, err := m.service.ParseToken(header)
		if err != nil {
			next.ServeHTTP(w, r)

			return
		}

		ctx := m.currentUserInformationProvider.SetCurrentUserID(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
