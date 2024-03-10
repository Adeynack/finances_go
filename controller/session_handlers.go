package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adeynack/finances/controller/routes"
	"github.com/adeynack/finances/model/query"
	view "github.com/adeynack/finances/view/templates/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ShowSession(c echo.Context) error {
	currentUser := useCurrentUser(c)
	if currentUser == nil {
		return Render(c, http.StatusOK, view.ShowLogin(routes.R.Session.Create.Path, view.ShowLoginOptions{}))
	}
	return Render(c, http.StatusOK, view.ShowSession(currentUser, routes.R.Session.Delete.Path))
}

func CreateSession(c echo.Context) error {
	db := useDb(c)
	u := query.Use(db).User

	email := c.FormValue("email")
	user, err := u.Where(u.Email.Eq(email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return invalidLogin(c, email)
	} else if err != nil {
		return fmt.Errorf("error loading user: %w", err)
	}

	password := c.FormValue("password")
	secret := useSecret(c)
	if checked, err := user.CheckPassword(secret, password); err != nil {
		return fmt.Errorf("error checking user's password: %w", err)
	} else if !checked {
		return invalidLogin(c, email)
	}

	session := useSession(c)
	session.Values[SessionCurrentUserId] = user.ID
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return c.Redirect(http.StatusFound, routes.R.Session.Show.Path)
}

func invalidLogin(c echo.Context, email string) error {
	return Render(c, http.StatusOK, view.ShowLogin(routes.R.Session.Create.Path, view.ShowLoginOptions{
		InvalidLogin: true,
		Email:        email,
	}))
}

func DeleteSession(c echo.Context) error {
	session := useSession(c)
	delete(session.Values, SessionCurrentUserId)
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return c.Redirect(http.StatusSeeOther, routes.R.Session.Show.Path)
}
