package interceptors

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicAccess = map[string]bool{
	"/grpc.BlogApi/Login":      true,
	"/grpc.BlogApi/CreateUser": true,
}

type AuthInterceptor struct {
	service                        service
	currentUserInformationProvider currentUserInformationProvider
}

func NewAuthInterceptor(s service, p currentUserInformationProvider) *AuthInterceptor {
	return &AuthInterceptor{
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

func (i *AuthInterceptor) UnaryAuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	if publicAccess[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is empty")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, errors.New("authorization token is not provided")
	}

	token := values[0]

	userID, err := i.service.ParseToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	ctx = i.currentUserInformationProvider.SetCurrentUserID(ctx, userID)

	return handler(ctx, req)
}
