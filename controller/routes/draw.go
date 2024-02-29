package routes

import (
	"net/http"

	"github.com/adeynack/finances/controller"
	"github.com/adeynack/finances/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type AppHandlerFunc func(ctx echo.Context, app *controller.RequestContext) error

func Draw(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Static("/view/static"))

	WithPainter(e.Group(""), forRequestContext(db), drawAppRoutes)
	WithPainter(e.Group(""), forAuthenticatedContext(db), drawAuthenticatedRoutes)
}

func drawAppRoutes(r RequestContextPainter[*controller.RequestContext]) {
	r.GET(PathRoot, controller.Root)
}

func drawAuthenticatedRoutes(r RequestContextPainter[*controller.AuthenticatedContext]) {

}

func createRequestContext(c echo.Context, db *gorm.DB) *controller.RequestContext {
	return &controller.RequestContext{
		DB: db.WithContext(c.Request().Context()),
	}
}

func forRequestContext(db *gorm.DB) HandlerGenerator[*controller.RequestContext] {
	return func(fn func(echo.Context, *controller.RequestContext) error) echo.HandlerFunc {
		return func(c echo.Context) error {
			return fn(c, createRequestContext(c, db))
		}
	}
}

func forAuthenticatedContext(db *gorm.DB) HandlerGenerator[*controller.AuthenticatedContext] {
	return func(fn func(echo.Context, *controller.AuthenticatedContext) error) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Truly implement. This is just garbage code to show it's possible to
			// 		   dynamically load a user from the session or request header.
			var user *model.User
			userId, err := c.Cookie("USER_ID")
			if err == nil {
				user = new(model.User)
				user.ID = userId.Value
			}
			if user == nil {
				return c.Redirect(http.StatusUnauthorized, "/need/to/authenticate")
			}

			return fn(c, &controller.AuthenticatedContext{
				RequestContext: createRequestContext(c, db),
			})
		}
	}
}
