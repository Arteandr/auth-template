package handlers

import (
	"mzhn/auth/internal/entity"
	mw "mzhn/auth/internal/middleware"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Logout(as *authservice.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := c.Get(mw.USER).(*entity.User)
		ctx := c.Request().Context()

		if err := as.Logout(ctx, user.Id); err != nil {
			return c.JSON(echo.ErrInternalServerError.Code, map[string]any{
				"message": err.Error(),
			})
		}

		return c.JSON(200, nil)
	}
}
