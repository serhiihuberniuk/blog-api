package interceptors

import (
	"context"
	"errors"
	"fmt"

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

type wrappedServerStream struct {
	grpc.ServerStream
	wrappedContext context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.wrappedContext
}

func wrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}

	return &wrappedServerStream{ServerStream: stream, wrappedContext: stream.Context()}
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

	userID, err := i.auth(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	ctx = i.currentUserInformationProvider.SetCurrentUserID(ctx, userID)

	return handler(ctx, req)
}

func (i *AuthInterceptor) StreamAuthInterceptor(srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	userID, err := i.auth(ss.Context())
	if err != nil {
		return status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	wrappedServerStream := wrapServerStream(ss)
	wrappedServerStream.wrappedContext = i.currentUserInformationProvider.
		SetCurrentUserID(wrappedServerStream.wrappedContext, userID)

	return handler(srv, wrappedServerStream)
}

func (i *AuthInterceptor) auth(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is empty")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", errors.New("authorization token is not provided")
	}

	token := values[0]

	userID, err := i.service.ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %w", err)
	}

	return userID, nil
}
