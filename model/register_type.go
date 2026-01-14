package model

import (
	"fmt"
	"strings"

	"github.com/adeynack/finances/app/appvalidator"
	"github.com/adeynack/finances/app/utils"
)

type RegisterType string

const (
	RegisterTypeIncome      = "Income"
	RegisterTypeExpense     = "Expense"
	RegisterTypeAsset       = "Asset"
	RegisterTypeBank        = "Bank"
	RegisterTypeCard        = "Card"
	RegisterTypeInstitution = "Institution"
	RegisterTypeInvestment  = "Investment"
	RegisterTypeLiability   = "Liability"
	RegisterTypeLoan        = "Loan"
)

var (
	CategoryTypes = []RegisterType{
		RegisterTypeIncome,
		RegisterTypeExpense,
	}
	AccountTypes = []RegisterType{
		RegisterTypeAsset,
		RegisterTypeBank,
		RegisterTypeCard,
		RegisterTypeInstitution,
		RegisterTypeInvestment,
		RegisterTypeLiability,
		RegisterTypeLoan,
	}
	RegisterTypes []RegisterType
)

func init() {
	RegisterTypes = append(RegisterTypes, CategoryTypes...)
	RegisterTypes = append(RegisterTypes, AccountTypes...)

	appvalidator.V.RegisterAlias("registerType", fmt.Sprintf("oneof=%s", strings.Join(utils.ToStringSlice(RegisterTypes), " ")))
}
