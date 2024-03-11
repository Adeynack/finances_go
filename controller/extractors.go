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
func useValueOf[T comparable](c echo.Context, key string, acceptNil bool) T {
	var defaultValue T
	rawValue := c.Get(key)
	if rawValue == nil {
		if acceptNil {
			return defaultValue
		} else {
			panic(fmt.Errorf("extracting %q from context: value not present", key))
		}
	}
	value, ok := rawValue.(T)
	if !ok {
		panic(fmt.Errorf("extracting %q from context: got %#v", key, value))
	}
	return value
}

// useAppContext returns the application's context object.
func useAppContext(c echo.Context) *AppContext {
	return useValueOf[*AppContext](c, CtxAppContext, false)
}

// useDb returns a GORM DB object in the context of the current requesst.
func useDb(c echo.Context) *gorm.DB {
	return useAppContext(c).DB.WithContext(c.Request().Context())
}

// useSecret returns the server's secret (eg: for encrypting users passwords).
func useSecret(c echo.Context) string {
	return useAppContext(c).ServerSecret
}

// useSession returns the actual session of the request.
func useSession(c echo.Context) *sessions.Session {
	session, err := session.Get(CtxSessionName, c)
	if err != nil {
		panic(fmt.Errorf("error getting session: %w", err))
	}
	return session
}

// useCurrentUser returns the currently authenticated user.
func useCurrentUser(c echo.Context) *model.User {
	return useValueOf[*model.User](c, CtxCurrentUser, true)
}
