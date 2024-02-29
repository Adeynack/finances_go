package controller

import (
	"os"

	"github.com/adeynack/finances/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RequestContext struct {
	DB     *gorm.DB
	Secret string
}

func NewRequestContext(ctx echo.Context, db *gorm.DB) RequestContext {
	return RequestContext{
		DB:     db.WithContext(ctx.Request().Context()),
		Secret: os.Getenv("SERVER_SECRET"),
	}
}

type AuthenticatedContext struct {
	*RequestContext
	CurrentUser *model.User
}

type BookContext struct {
	*AuthenticatedContext
	CurrentBook *model.Book
}
