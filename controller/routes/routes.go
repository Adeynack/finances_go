package routes

import "github.com/labstack/echo/v4"

type NamedRoutes struct {
	Root    *echo.Route
	Session struct {
		Show   *echo.Route
		Create *echo.Route
		Delete *echo.Route
	}
	Books struct {
		Index *echo.Route
	}
}

var R NamedRoutes
