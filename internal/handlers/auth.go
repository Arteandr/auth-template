package handlers

import (
	"errors"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/responses"
	"mzhn/auth/internal/middleware"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Auth(as *authservice.AuthService) echo.HandlerFunc {
	type request struct {
		Roles []entity.Role `json:"roles"`
	}

	return func(c echo.Context) error {
		var req request

		token := c.Get(middleware.TOKEN)
		if token != nil {
			return c.JSON(echo.ErrUnauthorized.Code, map[string]any{
				"error": "unathorized",
			})
		}

		if err := c.Bind(&req); err != nil {
			return err
		}

		ctx := c.Request().Context()

		_, err := as.Authenticate(ctx, &dto.Authenticate{
			AccessToken: token.(string),
			Roles:       req.Roles,
		})
		if err != nil {
			if errors.Is(err, authservice.ErrInsufficientPermission) {
				return responses.Forbidden(c)
			}

			return responses.Internal(c, err)
		}

		return responses.Ok(c, responses.Payload{})
	}
}
