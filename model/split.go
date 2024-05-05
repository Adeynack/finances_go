package model

type Split struct {
	BaseModel
	ExchangeID        string         `gorm:"type:uuid,not null" json:"exchange_id" validate:"required"`
	Exchange          *Exchange      `json:"exchange"`
	RegisterID        string         `gorm:"type:uuid,not null" json:"register_id" validate:"required"`
	Amount            int64          `gorm:"not null" json:"amount" validate:"required"`
	CounterpartAmount int64          `json:"counterpart_amount,omitempty"`
	Memo              string         `json:"memo,omitempty"`
	Status            ExchangeStatus `json:"status" validate:"required,exchangeStatus"`
}
