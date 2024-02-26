package model

import (
	"encoding/json"
	"fmt"
)

type BaseModel struct {
	ID        string      `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt ISODateTime `gorm:"not null" json:"created_at"`
	UpdatedAt ISODateTime `gorm:"not null" json:"updated_at"`
}

// MustJSON marhals the model to a JSON string. As it panics if an error occurs,
// it is meant as a quick debug convenience. Use proper `json.Marshal` with error
// handling for production-path JSON serialization.
func (model BaseModel) MustJSON() []byte {
	marshaled, err := json.Marshal(model)
	if err != nil {
		panic(fmt.Sprintf("error mashaling model %T to JSON: %v", model, err))
	}
	return marshaled
}
