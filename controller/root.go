package controller

import (
	"net/http"

	"github.com/adeynack/finances/view/templates/layouts"
	"github.com/labstack/echo/v4"
)

func Root(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, layouts.Page("Hello, templated web 2!"))
}
