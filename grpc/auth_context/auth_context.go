package auth_context

import (
	"context"
)

type (
	contextKey struct{}
)

var (
	contextSSOKey    = contextKey{}
	contextUserIDKey = contextKey{}
)

func WithSso(ctx context.Context, sso string) context.Context {
	ctx = context.WithValue(ctx, contextSSOKey, sso)
	return ctx
}

// SsoID returns the authenticated sso identifier associated with `ctx`.
func SsoID(ctx context.Context) string {
	if s, ok := ctx.Value(contextSSOKey).(string); ok {
		return s
	}
	return ""
}

func WithUserId(ctx context.Context, userId uint) context.Context {
	ctx = context.WithValue(ctx, contextUserIDKey, userId)
	return ctx
}

func UserId(ctx context.Context) uint {
	if s, ok := ctx.Value(contextUserIDKey).(uint); ok {
		return s
	}
	return 0
}
