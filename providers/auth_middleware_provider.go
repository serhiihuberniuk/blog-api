package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/serhiihuberniuk/blog-api/models"
)

const (
	authorizationHeader  = "Authorization"
	bearerAuthentication = "bearer"
)

type AuthMiddlewareProvider struct {
	service                        service
	currentUserInformationProvider currentUserInformationProvider
}

func NewAuthInfoProvider(s service, p currentUserInformationProvider) *AuthMiddlewareProvider {
	return &AuthMiddlewareProvider{
		service:                        s,
		currentUserInformationProvider: p,
	}
}

type service interface {
	ParseToken(tokenString string) (string, error)
}

type currentUserInformationProvider interface {
	SetCurrentUserID(ctx context.Context, userID string) context.Context
}

func (p *AuthMiddlewareProvider) BearerAuthMiddleware(r *http.Request) (context.Context, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return r.Context(), models.ErrNotAuthenticated
	}

	headerSplit := strings.Split(header, " ")
	if len(headerSplit) != 2 {
		return r.Context(), models.ErrNotAuthenticated
	}

	if strings.ToLower(headerSplit[0]) != bearerAuthentication {
		return r.Context(), models.ErrNotAuthenticated
	}

	userID, err := p.service.ParseToken(headerSplit[1])
	if err != nil {
		return r.Context(), fmt.Errorf("error occurred while parsing token: %w", err)
	}

	ctx := p.currentUserInformationProvider.SetCurrentUserID(r.Context(), userID)

	return ctx, nil
}
