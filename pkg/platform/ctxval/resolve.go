package ctxval

import (
	"context"
	"fmt"
)

// Resolve finds in the provided ctx for an unnamed resource of type T.
func Resolve[T any](ctx context.Context) (T, error) {
	return ResolveNamed[T](ctx, "")
}

func ResolveNamed[T any](ctx context.Context, name string) (T, error) {
	var result T

	key := keyFor[T](name)

	rawValue := ctx.Value(key)
	if rawValue == nil {
		return result, fmt.Errorf("%w %q", ErrUnregisteredDependency, key)
	}

	result, ok := rawValue.(T)
	if !ok {
		return result, fmt.Errorf("%w %q", ErrUnexpectedType, key)
	}

	return result, nil
}

func MustResolve[T any](ctx context.Context) T {
	return MustResolveNamed[T](ctx, "")
}

func MustResolveNamed[T any](ctx context.Context, name string) T {
	result, err := ResolveNamed[T](ctx, name)
	if err != nil {
		panic(err)
	}

	return result
}
