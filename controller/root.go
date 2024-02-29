package controller

import (
	"fmt"
	"net/http"

	"github.com/adeynack/finances/model/query"
	"github.com/adeynack/finances/view/templates/layouts"
	"github.com/labstack/echo/v4"
)

func Root(ctx echo.Context, app *RequestContext) error {
	u := query.Use(app.DB).User
	usersWithBooks, err := u.Preload(u.Books).Find()
	if err != nil {
		return fmt.Errorf("error loading users with books: %v", err)
	}
	return Render(ctx, http.StatusOK, layouts.Page("Hello, templated web 2!", usersWithBooks))
}
