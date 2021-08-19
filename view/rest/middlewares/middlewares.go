package middlewares

import "context"

type Middleware struct {
	service
	contextValueProvider
}

func NewMiddleware(s service, p contextValueProvider) *Middleware {
	return &Middleware{
		service:              s,
		contextValueProvider: p,
	}
}

type contextValueProvider interface {
	SetCurrentUserID(ctx context.Context, userID string) context.Context
}

type service interface {
	ParseToken(tokenString string) (string, error)
}
