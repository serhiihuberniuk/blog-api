package graphqlMiddlewares

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

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.authMiddlewareProvider.BearerAuthMiddleware(r)
		if err != nil {
			next.ServeHTTP(w, r)

			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
