package model

import (
	"fmt"
	"reflect"
)

// All lists all models that are to be managed by Gorm
func All() []any {
	return []any{
		ImportOrigin{},
		User{},
		Book{},
		Register{},
		Exchange{},
	}
}

func ModelAndID(subject any) (model string, id string, err error) {
	v := reflect.ValueOf(subject)
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("Value of type %T is not a struct", subject)
		return
	}

	model = v.Type().Name()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		err = fmt.Errorf("Value of type %T does not contain field \"ID\"", subject)
		return
	}

	id = idField.String()
	return
}
