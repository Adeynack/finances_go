package model

import (
	"fmt"
	"strings"
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
