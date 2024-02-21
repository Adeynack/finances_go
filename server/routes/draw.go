package routes

import (
	"github.com/adeynack/finances/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Draw(e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Static("/static"))

	e.GET(PathRoot, handlers.Root)
}
