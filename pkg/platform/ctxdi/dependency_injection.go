package ctxdi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type dependencyKey struct {
	typ  reflect.Type
	name string
}

// Errors
var (
	ErrCannotResolve          = errors.New("unable to resolve dependency")
	ErrUnregisteredDependency = fmt.Errorf("%w: unregistered dependency", ErrCannotResolve)
	ErrUnexpectedType         = fmt.Errorf("%w: unexpected type from registered dependency", ErrCannotResolve)
	ErrInvalidName            = fmt.Errorf("%w: invalid name", ErrCannotResolve)
)

func keyFor[T any](name string) dependencyKey {
	return dependencyKey{
		typ:  reflect.TypeFor[T](),
		name: name,
	}
}

func (k dependencyKey) String() string {
	typePart := k.typ.String()

	if k.name == "" {
		return typePart
	}

	return fmt.Sprintf("%s(%s)", typePart, k.name)
}

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

func RegisterValue[T any](ctx context.Context, value T) context.Context {
	return RegisterNamedValue(ctx, "", value)
}

func RegisterNamedValue[T any](ctx context.Context, name string, value T) context.Context {
	key := keyFor[T](name)
	return context.WithValue(ctx, key, value)
}
