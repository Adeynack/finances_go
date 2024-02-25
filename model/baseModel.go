package model

type baseModel struct {
	ID        string      `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt ISODateTime `gorm:"not null" json:"created_at"`
	UpdatedAt ISODateTime `gorm:"not null" json:"updated_at"`
}
