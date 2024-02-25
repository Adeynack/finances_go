package db

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TrappedMigrationsError struct {
	PendingMigrations []string
}

// Error implements error.
func (*TrappedMigrationsError) Error() string {
	return "pending migrations detected" // TODO: Instruct to use `cmd/tool` (still to be coded) to output the missing SQL migrations
}

type MigrationTrapDialector struct {
	gorm.Dialector
	migrator *MigrationTrapMigrator
}

func NewMigrationTrapDialector(baseDialector gorm.Dialector) *MigrationTrapDialector {
	return &MigrationTrapDialector{Dialector: baseDialector}
}

func (d *MigrationTrapDialector) Migrator(db *gorm.DB) gorm.Migrator {
	if d.migrator == nil {
		d.migrator = &MigrationTrapMigrator{Migrator: d.Dialector.Migrator(db)}
	}
	return d.migrator
}

type MigrationTrapMigrator struct {
	gorm.Migrator
	trappedMigrations []string
}

func (m *MigrationTrapMigrator) addChange(format string, args ...any) error {
	m.trappedMigrations = append(m.trappedMigrations, fmt.Sprintf(format, args...))
	return nil
}

func (m *MigrationTrapMigrator) AutoMigrate(dst ...interface{}) error {
	err := m.Migrator.AutoMigrate(dst...)
	if err != nil {
		return err
	}
	if len(m.trappedMigrations) > 0 {
		return &TrappedMigrationsError{PendingMigrations: m.trappedMigrations}
	}
	return nil
}

// AddColumn implements gorm.Migrator.
func (m *MigrationTrapMigrator) AddColumn(dst interface{}, field string) error {
	return m.addChange("Add column: %v %v", dst, field)
}

// AlterColumn implements gorm.Migrator.
func (m *MigrationTrapMigrator) AlterColumn(dst interface{}, field string) error {
	return m.addChange("Alter column: %v %v", dst, field)
}

// CreateConstraint implements gorm.Migrator.
func (m *MigrationTrapMigrator) CreateConstraint(dst interface{}, name string) error {
	return m.addChange("Create constraint: %v %v", dst, name)
}

// CreateIndex implements gorm.Migrator.
func (m *MigrationTrapMigrator) CreateIndex(dst interface{}, name string) error {
	return m.addChange("Create index: %v, %v", dst, name)
}

// CreateTable implements gorm.Migrator.
func (m *MigrationTrapMigrator) CreateTable(dst ...interface{}) error {
	return m.addChange("Create table: %s", dstTypeNames(dst))
}

// CreateView implements gorm.Migrator.
func (m *MigrationTrapMigrator) CreateView(name string, option gorm.ViewOption) error {
	return m.addChange("Create view: %v %v", name, option)
}

// DropColumn implements gorm.Migrator.
func (m *MigrationTrapMigrator) DropColumn(dst interface{}, field string) error {
	return m.addChange("Drop column: %v %v", dst, field)
}

// DropConstraint implements gorm.Migrator.
func (m *MigrationTrapMigrator) DropConstraint(dst interface{}, name string) error {
	return m.addChange("Drop constraint: %v %v", dst, name)
}

// DropIndex implements gorm.Migrator.
func (m *MigrationTrapMigrator) DropIndex(dst interface{}, name string) error {
	return m.addChange("Drop index: %v, %v", dst, name)
}

// DropTable implements gorm.Migrator.
func (m *MigrationTrapMigrator) DropTable(dst ...interface{}) error {
	return m.addChange("Drop table: %s", dstTypeNames(dst))
}

// DropView implements gorm.Migrator.
func (m *MigrationTrapMigrator) DropView(name string) error {
	return m.addChange("Drop view: %v", name)
}

// MigrateColumn implements gorm.Migrator.
func (m *MigrationTrapMigrator) MigrateColumn(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	return m.addChange("Migrate column: %T %v %v", dst, field, columnType)
}

// MigrateColumnUnique implements gorm.Migrator.
func (m *MigrationTrapMigrator) MigrateColumnUnique(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	return m.addChange("Migrate column unique: %v %v %v", dst, field, columnType)
}

// RenameColumn implements gorm.Migrator.
func (m *MigrationTrapMigrator) RenameColumn(dst interface{}, oldName string, field string) error {
	return m.addChange("Rename column: %v %v %v", dst, oldName, field)
}

// RenameIndex implements gorm.Migrator.
func (m *MigrationTrapMigrator) RenameIndex(dst interface{}, oldName string, newName string) error {
	return m.addChange("Rename index: %v %v %v", dst, oldName, newName)
}

// RenameTable implements gorm.Migrator.
func (m *MigrationTrapMigrator) RenameTable(oldName interface{}, newName interface{}) error {
	return m.addChange("Rename table: %v %v", oldName, newName)
}

func dstTypeNames(dst []any) string {
	typeNames := make([]string, len(dst))
	for index, element := range dst {
		typeNames[index] = fmt.Sprintf("%T", element)
	}
	return strings.Join(typeNames, " ")
}
