package model

import (
	"time"
)

type Exchange struct {
	BaseModel
	Date        time.Time      `gorm:"not null,type:date" json:"date" validate:"required"`
	RegisterID  string         `gorm:"type:uuid,not null" json:"register_id" validate:"required"`
	Register    *Register      `json:"register,omitempty"`
	Cheque      string         `json:"cheque,omitempty"`
	Description string         `gorm:"not null" json:"description" validate:"required"`
	Memo        string         `json:"memo,omitempty"`
	Status      ExchangeStatus `gorm:"not null" json:"status" validate:"required,exchangeStatus"`
}
