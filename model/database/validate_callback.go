package database

import (
	"reflect"

	"github.com/adeynack/finances/app/appvalidator"
	"gorm.io/gorm"
)

func validateCallback(db *gorm.DB) {
	dest := db.Statement.Dest
	v := reflect.ValueOf(dest)
	for i := 0; i < v.Len(); i++ {
		model := v.Index(i).Interface()
		err := appvalidator.V.StructCtx(db.Statement.Context, model)
		if err != nil {
			_ = db.AddError(err)
		}
	}
}
