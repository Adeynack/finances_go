package utils

import (
	"fmt"

	"gorm.io/gorm"
)

func GormTableForModel(db *gorm.DB, model any) (string, error) {
	// Credits: https://stackoverflow.com/a/64385175
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(model)
	if err != nil {
		return "", fmt.Errorf("error parsing model: %w", err)
	}
	return stmt.Schema.Table, nil
}

func GormTablesForModels(db *gorm.DB, models ...any) ([]string, error) {
	tableNames := make([]string, len(models))
	for index, model := range models {
		table, err := GormTableForModel(db, model)
		if err != nil {
			return nil, err
		}
		tableNames[index] = table
	}
	return tableNames, nil
}
