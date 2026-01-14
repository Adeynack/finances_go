package controller

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Render outputs a component (compatible with the templ.Component
// interface) to an Echo response.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
