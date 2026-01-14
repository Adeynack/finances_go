package controller

import (
	"net/http"

	"github.com/adeynack/finances/view/templates/layouts"
	"github.com/labstack/echo/v4"
)

func Root(c echo.Context) error {
	return Render(c, http.StatusOK, layouts.Page())
}
