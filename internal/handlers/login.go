package handlers

import (
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Login(as *authservice.AuthService) echo.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	return func(c echo.Context) error {
		var req request

		if err := c.Bind(&req); err != nil {
			return err
		}

		tokens, err := as.Login(c.Request().Context(), &dto.Login{
			Email:    req.Email,
			Password: req.Password,
		})
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
