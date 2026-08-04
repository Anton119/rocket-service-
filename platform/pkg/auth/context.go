package auth

import "context"

type userUUIDCtxKey struct{}

type sessionUUIDCtxKey struct{}

// WithUserUUID кладёт UUID пользователя в контекст.
func WithUserUUID(ctx context.Context, userUUID string) context.Context {
	return context.WithValue(ctx, userUUIDCtxKey{}, userUUID)
}

// UserUUIDFromContext извлекает UUID пользователя из контекста.
func UserUUIDFromContext(ctx context.Context) (string, bool) {
	userUUID, ok := ctx.Value(userUUIDCtxKey{}).(string)
	return userUUID, ok
}

// WithSessionUUID кладёт UUID сессии в контекст.
func WithSessionUUID(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionUUIDCtxKey{}, sessionUUID)
}

// SessionUUIDFromContext извлекает UUID сессии из контекста.
func SessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDCtxKey{}).(string)
	return sessionUUID, ok
}
