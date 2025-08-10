package middleware

import "context"

// 使用类型化上下文键，避免字符串键冲突
type contextKey string

const (
	CtxKeyRequestID contextKey = "requestId"
	CtxKeyUser      contextKey = "user"
	CtxKeyUserID    contextKey = "userId"
	CtxKeyUsername  contextKey = "username"
	CtxKeyRoles     contextKey = "roles"
	CtxKeyPerms     contextKey = "permissions"
	CtxKeyToken     contextKey = "token"
)

func WithValue(ctx context.Context, key contextKey, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

func Value[T any](ctx context.Context, key contextKey) (T, bool) {
	v, ok := ctx.Value(key).(T)
	var zero T
	if !ok {
		return zero, false
	}
	return v, true
}
