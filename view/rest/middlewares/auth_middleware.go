package middlewares

import (
	"context"
	"net/http"
)

type AuthMiddleware struct {
	authMiddlewareProvider authMiddlewareProvider
}

func NewAuthMiddleware(p authMiddlewareProvider) *AuthMiddleware {
	return &AuthMiddleware{
		authMiddlewareProvider: p,
	}
}

type authMiddlewareProvider interface {
	BearerAuthMiddleware(r *http.Request) (context.Context, error)
}

func (m *AuthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.authMiddlewareProvider.BearerAuthMiddleware(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next(w, r.WithContext(ctx))
	}
}
