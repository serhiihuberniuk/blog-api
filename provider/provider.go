package provider

import (
	"context"
	"fmt"
)

const userIdKey contextKey = "userId"

type contextKey string

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) SetCurrentUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIdKey, userID)
}

func (p *Provider) GetCurrentUserID(ctx context.Context) string {
	return fmt.Sprint(ctx.Value(userIdKey))
}
