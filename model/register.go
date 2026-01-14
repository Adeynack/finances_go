package model

import "gorm.io/datatypes"

type Register struct {
	BaseModel
	Name               string          `gorm:"not null" json:"name" validate:"required"`
	Type               string          `gorm:"not null,enum:" json:"type" validate:"required,registerType"`
	BookID             string          `gorm:"type:uuid,not null,index" json:"book_id" validate:"required"`
	Book               *Book           `json:"omitempty"`
	ParentID           string          `gorm:"type:uuid" json:"parent_id"`
	Parent             *Register       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	StartsAt           datatypes.Date  `gorm:"not null" json:"starts_at" validate:"required"`
	ExpiresAt          *datatypes.Date `json:"expires_at"`
	CurrentcyIsoCode   ISOCurrency     `gorm:"not null" json:"currency_iso_code" validate:"currencyCode"`
	Notes              string          `json:"notes"`
	InitialBalance     int64           `gorm:"not null" json:"initial_balance"`
	Active             bool            `gorm:"not null" json:"active"`
	DefaultCategoryID  string          `gorm:"indexed" json:"default_category_id"`
	DefaultCategory    *Register       `gorm:"foreignKey:DefaultCategoryID" json:"default_category,omitempty"`
	InstitutionName    string          `json:"institution_name,omitempty"`
	AccountNumber      string          `json:"account_number,omitempty"`
	AnnualInterestRate float32         `json:"annual_interest_rate,omitempty"`
	CreditLimit        int64           `json:"credit_limit,omitempty"`
	CardNumber         string          `json:"card_number,omitempty"`
	ImportOrigins      []*ImportOrigin `gorm:"polymorphic:Subject" json:"import_origins,omitempty"`
}
