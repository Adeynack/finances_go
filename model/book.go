package model

type Book struct {
	BaseModel
	Name                   string          `gorm:"not null" json:"name" validate:"required"`
	OwnerID                string          `gorm:"not null,index" json:"owner_id" validate:"required"`
	Owner                  *User           `json:"owner,omitempty"`
	DefaultCurrencyIsoCode ISOCurrency     `gorm:"not null,type:char(3)" json:"default_currency_iso_code" validate:"currencyCode"`
	ImportOrigins          []*ImportOrigin `gorm:"polymorphic:Subject" json:"import_origins"`
}
