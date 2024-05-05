package model

import (
	"fmt"
	"strings"

	"github.com/adeynack/finances/app/appvalidator"
)

type ISOCurrency string

var (
	_ fmt.Stringer = (*ISOCurrency)(nil)
)

// String implements fmt.Stringer.
func (i ISOCurrency) String() string {
	// Ensures display and saving to the DB is done in upper case.
	return strings.ToUpper(string(i))
}

func init() {
	appvalidator.V.RegisterAlias("currencyCode", "alpha,uppercase,len=3")
}
