package ctxdi

import (
	"fmt"
	"reflect"
)

type dependencyKey struct {
	typ  reflect.Type
	name string
}

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
