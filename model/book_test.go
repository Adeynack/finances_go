package model

import (
	"fmt"
	"testing"

	"github.com/adeynack/finances/app/appvalidator"
	"github.com/stretchr/testify/assert"
)

func Test_Book_ValidateISOCurrency_LowerCase(t *testing.T) {
	book := Book{DefaultCurrencyIsoCode: "EU"}
	validations, err := appvalidator.V.Namespaced(appvalidator.V.Struct(book))
	if !assert.NoError(t, err) {
		return
	}
	validation, ok := validations["Book.DefaultCurrencyIsoCode"]
	if !assert.Equal(t, true, ok) {
		return
	}
	fmt.Println(validation)
}
