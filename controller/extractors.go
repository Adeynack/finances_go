package controller

import (
	"fmt"

	"github.com/adeynack/finances/model"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// useValueOf extract the value for provided key from the request context
// and expects it to be of type T. It panics if the value was not found or
// it could not be cast to T, with proper message.
//
// It uses panic as an error mechanism since it's not even supposed to at
// this point. Those values are expected to be there and if they do not,
// it's detected and fixed during development and testing. It's not on the
// production's hot-path, and therefore acceptable to lightened the call-site
// by skipping error testing.
func useValueOf[T any](c echo.Context, key string) T {
	value, ok := c.Get(key).(T)
	if !ok {
		panic(fmt.Errorf("extracting %q from context: got %#v", key, value))
	}
	return value
}

func useDb(c echo.Context) *gorm.DB {
	return useValueOf[*gorm.DB](c, CtxRequestDb)
}

func useSecret(c echo.Context) string {
	return useValueOf[string](c, CtxServerSecret)
}

func useSession(c echo.Context) *sessions.Session {
	session, err := session.Get(CtxSessionName, c)
	if err != nil {
		panic(fmt.Errorf("error getting session: %v", err))
	}
	return session
}

func useCurrentUser(c echo.Context) *model.User {
	value := c.Get(CtxCurrentUser)
	if value == nil {
		return nil
	}
	currentUser, ok := value.(*model.User)
	if !ok {
		panic(fmt.Errorf("extracting %q from context: got %#v", CtxCurrentUser, value))
	}
	return currentUser
}
