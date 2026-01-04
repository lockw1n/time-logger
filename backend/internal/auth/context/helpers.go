package context

import "context"

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, auth)
}

func FromContext(ctx context.Context) (AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey).(AuthContext)
	return auth, ok
}
