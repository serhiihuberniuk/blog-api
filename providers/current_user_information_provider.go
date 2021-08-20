package providers

import "context"

const userIdKey contextKey = "userId"

type contextKey string

type CurrentUserInformationProvider struct{}

func NewCurrentUserInformationProvider() *CurrentUserInformationProvider {
	return &CurrentUserInformationProvider{}
}

func (p *CurrentUserInformationProvider) SetCurrentUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIdKey, userID)
}

func (p *CurrentUserInformationProvider) GetCurrentUserID(ctx context.Context) string {
	return ctx.Value(userIdKey).(string)
}
