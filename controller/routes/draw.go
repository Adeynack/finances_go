package routes

import (
	"github.com/adeynack/finances/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type AppHandlerFunc func(ctx echo.Context, app *controller.RequestContext) error

func Draw(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Static("/view/static"))

	inContext := func(fn AppHandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return fn(c, &controller.RequestContext{
				DB: db.WithContext(c.Request().Context()),
			})
		}
	}

	e.GET(PathRoot, inContext(controller.Root))
}
