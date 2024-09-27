package handlers

import (
	"mzhn/auth/internal/dto"
	mw "mzhn/auth/internal/middleware"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Refresh(as *authservice.AuthService) echo.HandlerFunc {

	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	return func(c echo.Context) error {

		token := c.Get(mw.TOKEN)
		if token == nil {
			return c.JSON(echo.ErrInternalServerError.Code, map[string]any{
				"message": "token not found",
			})
		}

		tokens, err := as.Refresh(c.Request().Context(), &dto.Refresh{RefreshToken: token.(string)})
		if err != nil {
			return c.JSON(echo.ErrInternalServerError.Code, map[string]any{
				"message": err.Error(),
			})
		}

		return c.JSON(200, &response{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
