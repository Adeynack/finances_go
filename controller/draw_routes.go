package controller

import (
	"path/filepath"

	"github.com/adeynack/finances/controller/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func DrawRoutes(e *echo.Echo, appContext *AppContext) {
	e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Static(filepath.Join("view", "static")),
	)

	routes.R.Root = e.GET("", Root)

	drawSessionRoutes(e.Group(""), appContext)
}

func drawSessionRoutes(e *echo.Group, appContext *AppContext) {
	e.Use(
		SetAppContext(appContext),
		session.Middleware(sessions.NewCookieStore([]byte(appContext.ServerSecret))),
		SetCurrentUser(),
	)

	routes.R.Session.Show = e.GET(routes.R.Root.Path+"/session", ShowSession)
	routes.R.Session.Create = e.POST(routes.R.Session.Show.Path, CreateSession)
	routes.R.Session.Delete = e.DELETE(routes.R.Session.Show.Path, DeleteSession)
}
