package controller

import (
	"errors"
	"fmt"

	"github.com/adeynack/finances/model/query"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	CtxCurrentUser       = "CurrentUser"
	CtxSessionName       = "Session"
	SessionCurrentUserId = "CurrentUserId"
)

func SetCurrentUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := useSession(c)
			currentUserId, ok := session.Values[SessionCurrentUserId].(string)
			if ok && len(currentUserId) > 0 {
				db := useDb(c)
				u := query.Use(db).User
				currentUser, err := u.Where(u.ID.Eq(currentUserId)).First()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// noop: User simply does not exists
				} else if err != nil {
					return fmt.Errorf("error loading current user: %w", err)
				}
				c.Set(CtxCurrentUser, currentUser)
			}

			return next(c)
		}
	}
}
