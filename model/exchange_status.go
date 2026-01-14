package model

import (
	"fmt"
	"strings"

	"github.com/adeynack/finances/app/appvalidator"
	"github.com/adeynack/finances/app/utils"
)

type ExchangeStatus string

const (
	ExchangeStatusUncleared   = "Uncleared"
	ExchangeStatusReconciling = "Reconciling"
	ExchangeStatusCleared     = "Cleared"
)

var (
	ExchangeStatuses = []ExchangeStatus{
		ExchangeStatusUncleared,
		ExchangeStatusReconciling,
		ExchangeStatusCleared,
	}
)

func init() {
	appvalidator.V.RegisterAlias("exchangeStatus", fmt.Sprintf("oneof=%s", strings.Join(utils.ToStringSlice(ExchangeStatuses), " ")))
}
