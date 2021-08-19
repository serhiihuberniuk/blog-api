package providers

import "context"

const userIdKey contextKey = "userId"

type contextKey string

type ContextValueProvider struct{}

func NewContextValueProvider() *ContextValueProvider {
	return &ContextValueProvider{}
}

func (p *ContextValueProvider) SetCurrentUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIdKey, userID)
}

func (p *ContextValueProvider) GetCurrentUserID(ctx context.Context) string {
	return ctx.Value(userIdKey).(string)
}
