package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	CtxServerSecret = "CTX_SERVER_SECRET"
	CtxRequestDb    = "CTX_REQUEST_DB"
)

func SetRequestContext(db *gorm.DB, serverSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set server-static values
			c.Set(CtxServerSecret, serverSecret)

			// Set values proper to the request.
			c.Set(CtxRequestDb, db.WithContext(c.Request().Context()))

			return next(c)
		}
	}
}
