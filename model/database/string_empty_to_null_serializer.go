package database

import (
	"context"
	"reflect"

	"gorm.io/gorm/schema"
)

type StringEmptyToNullSerializer struct{}

var _ schema.SerializerInterface = StringEmptyToNullSerializer{}

// Scan implements schema.SerializerInterface.
func (s StringEmptyToNullSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	panic("unimplemented")
}

// Value implements schema.SerializerInterface.
func (s StringEmptyToNullSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	panic("unimplemented")
}
