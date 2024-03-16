package appvalidator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AppValidator struct {
	*validator.Validate
}

var V *AppValidator

func init() {
	V = &AppValidator{Validate: validator.New()}
	V.RegisterAlias("currencyCode", "alpha,uppercase,len=3")
}

// Namespaced returns a map of namespaced FieldError if the provided err
// was validator.ValidationErrors. If it was not, returns an empty map,
// but then the error as second return value.
func (v *AppValidator) Namespaced(err error) (map[string]validator.FieldError, error) {
	if err == nil {
		return nil, nil
	}

	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		return nil, err
	}

	errorPerField := make(map[string]validator.FieldError)
	for _, fieldError := range validationErrors {
		if _, exists := errorPerField[fieldError.Namespace()]; exists {
			// Just for the time of development, until checked if 2 errors on the same field have 2 different errors in the slice.
			panic(fmt.Errorf("namespace %q already has an error", fieldError.Namespace()))
		}
		errorPerField[fieldError.Namespace()] = fieldError
	}
	return errorPerField, nil
}
