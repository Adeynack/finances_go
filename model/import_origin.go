package model

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type ImportOrigin struct {
	BaseModel
	SubjectType    string `gorm:"not null" json:"subject_type"`
	SubjectID      string `gorm:"not null" json:"subject_id"`
	ExternalSystem string `gorm:"not null" json:"external_system"`
	ExternalID     string `gorm:"not null" json:"external_id"`
}

func (o *ImportOrigin) SetSubject(subject any) error {
	subjectReflValue := reflect.ValueOf(subject)
	subjectReflType := subjectReflValue.Type()

	if field, ok := subjectReflType.FieldByName("ID"); ok {
		o.SubjectID = subjectReflValue.FieldByIndex(field.Index).String()
	} else {
		return fmt.Errorf("subject of type %T does not have an ID field", subject)
	}

	o.SubjectType = subjectReflType.Name()

	return nil
}

func CreateImportOrigin(db *gorm.DB, subject any, externalSystem string, externalID string) error {
	model, id, err := ModelAndID(subject)
	if err != nil {
		return fmt.Errorf("error creating new ImportOrigin: %w", err)
	}

	origin := &ImportOrigin{
		SubjectType:    model,
		SubjectID:      id,
		ExternalSystem: externalSystem,
		ExternalID:     externalID,
	}

	err = db.Create(origin).Error
	if err != nil {
		return fmt.Errorf("error creating new ImportOrigin: %w", err)
	}

	return err
}
