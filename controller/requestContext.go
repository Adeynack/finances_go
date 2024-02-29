package controller

import (
	"gorm.io/gorm"
)

type RequestContext struct {
	DB *gorm.DB
}
