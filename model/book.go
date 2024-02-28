package model

type Book struct {
	BaseModel
	Name                   string `gorm:"not null" json:"name"`
	OwnerID                string `gorm:"not null,index" json:"owner_id"`
	Owner                  *User  `json:"owner,omitempty"`
	DefaultCurrencyIsoCode string `gorm:"not null,type:char(3)" json:"default_currency_iso_code"`
}
