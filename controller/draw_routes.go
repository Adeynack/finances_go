package controller

import (
	"path/filepath"

	"github.com/adeynack/finances/controller/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func DrawRoutes(e *echo.Echo, db *gorm.DB, serverSecret string) {
	e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Static(filepath.Join("view", "static")),
	)

	routes.R.Root = e.GET("", Root)

	drawSessionRoutes(e.Group(""), db, serverSecret)
}

func drawSessionRoutes(e *echo.Group, db *gorm.DB, serverSecret string) {
	e.Use(
		SetRequestContext(db, serverSecret),
		session.Middleware(sessions.NewCookieStore([]byte(serverSecret))),
		SetCurrentUser(),
	)

	routes.R.Session.Show = e.GET(routes.R.Root.Path+"/session", ShowSession)
	routes.R.Session.Create = e.POST(routes.R.Session.Show.Path, CreateSession)
	routes.R.Session.Delete = e.DELETE(routes.R.Session.Show.Path, DeleteSession)
}
