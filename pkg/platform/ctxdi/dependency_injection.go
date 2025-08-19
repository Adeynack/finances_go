package ctxdi

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrCannotResolve          = errors.New("unable to resolve dependency")
	ErrUnregisteredDependency = fmt.Errorf("%w: unregistered dependency", ErrCannotResolve)
	ErrUnexpectedType         = fmt.Errorf("%w: unexpected type from registered dependency", ErrCannotResolve)
	ErrInvalidName            = fmt.Errorf("%w: invalid name", ErrCannotResolve)
)
