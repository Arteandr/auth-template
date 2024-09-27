package handlers

import (
	"log/slog"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
	mw "mzhn/auth/internal/middleware"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Profile(as *authservice.AuthService) echo.HandlerFunc {

	type response struct {
		Id         string        `json:"id"`
		LastName   *string       `json:"lastName"`
		FirstName  *string       `json:"firstName"`
		MiddleName *string       `json:"middleName"`
		Email      string        `json:"email"`
		Roles      []entity.Role `json:"roles"`
	}

	return func(c echo.Context) error {
		claims := c.Get(mw.USER).(*entity.User)
		ctx := c.Request().Context()

		user, roles, err := as.Profile(ctx, claims.Id)
		if err != nil {
			slog.Error("cannot profile user", sl.Err(err))
			return c.JSON(echo.ErrInternalServerError.Code, map[string]any{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &response{
			Id:         user.Id,
			LastName:   user.LastName,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			Email:      user.Email,
			Roles:      roles,
		})
	}
}
