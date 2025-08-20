package ctxval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar string
}

type Bar struct {
	Foo string
}

func TestResolve(t *testing.T) {
	t.Run("when nothing is registered", func(t *testing.T) {
		t.Run("Resolve fails", func(t *testing.T) {
			_, err := Resolve[Foo](t.Context())
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("ResolveNamed fails", func(t *testing.T) {
			_, err := ResolveNamed[Foo](t.Context(), "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})

		t.Run("MustResolve panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo"`,
				func() { MustResolve[Foo](t.Context()) },
			)
		})

		t.Run("MustResolveNamed panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`,
				func() { MustResolveNamed[Foo](t.Context(), "asdf") },
			)
		})
	})

	t.Run("when an unnamed Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "5160b303-f563-44c3-ac93-baebea18cbe7"}
		ctx := RegisterValue(t.Context(), registeredFoo)

		t.Run("Resolve succeeds for unnamed Foo", func(t *testing.T) {
			result, err := Resolve[Foo](ctx)
			require.NoError(t, err)
			require.Equal(t, registeredFoo, result)
		})

		t.Run("Resolve fails for unnamed Bar", func(t *testing.T) {
			_, err := Resolve[Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Bar"`, err.Error())
		})

		t.Run("ResolveNamed fails for named Foo", func(t *testing.T) {
			_, err := ResolveNamed[Foo](ctx, "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})
	})

	t.Run("when a named Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "e1950227-441b-4238-804f-908110c0592a"}
		ctx := RegisterNamedValue(t.Context(), "TheFuu", registeredFoo)

		t.Run("Resolve fails for unnamed Foo", func(t *testing.T) {
			_, err := Resolve[Foo](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("Resolve fails for unnamed Bar", func(t *testing.T) {
			_, err := Resolve[Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Bar"`, err.Error())
		})

		t.Run("Resolve succeeds for Foo named 'TheFuu'", func(t *testing.T) {
			result, err := ResolveNamed[Foo](ctx, "TheFuu")
			require.NoError(t, err)
			require.Equal(t, registeredFoo, result)
		})

		t.Run("Resolve fails for Foo named 'SomethingElse'", func(t *testing.T) {
			_, err := ResolveNamed[Foo](ctx, "SomethingElse")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(SomethingElse)"`, err.Error())
		})
	})
}
