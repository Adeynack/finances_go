package ctxdi

import "context"

func RegisterValue[T any](ctx context.Context, value T) context.Context {
	return RegisterNamedValue(ctx, "", value)
}

func RegisterNamedValue[T any](ctx context.Context, name string, value T) context.Context {
	key := keyFor[T](name)
	return context.WithValue(ctx, key, value)
}
