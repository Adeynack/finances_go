package model

import (
	"testing"

	"github.com/adeynack/finances/app/appvalidator"
	"github.com/stretchr/testify/require"
)

func Test_Book_ValidateISOCurrency_LowerCase(t *testing.T) {
	book := Book{DefaultCurrencyIsoCode: "EU"}

	validations, err := appvalidator.V.Namespaced(appvalidator.V.Struct(book))
	require.NoError(t, err)

	validationError, ok := validations["Book.DefaultCurrencyIsoCode"]
	require.True(t, ok)
	require.Equal(t, validationError.ActualTag(), "len")
	require.Equal(t, validationError.Param(), "3")
	require.Equal(t, validationError.Field(), "DefaultCurrencyIsoCode")
}
