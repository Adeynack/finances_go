package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Root(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, Web!")
}
