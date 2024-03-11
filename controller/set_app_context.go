package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	CtxAppContext = "AppContext"
)

type AppContext struct {
	ServerSecret string
	DB           *gorm.DB
}

func SetAppContext(appContext *AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(CtxAppContext, appContext)
			return next(c)
		}
	}
}
