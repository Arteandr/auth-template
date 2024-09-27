package handlers

import (
	"errors"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/responses"
	"mzhn/auth/internal/services/authservice"

	"github.com/labstack/echo/v4"
)

func Register(as *authservice.AuthService) echo.HandlerFunc {
	type request struct {
		LastName   *string       `json:"lastName"`
		FirstName  *string       `json:"firstName"`
		MiddleName *string       `json:"middleName"`
		Email      string        `json:"email"`
		Password   string        `json:"password"`
		Roles      []entity.Role `json:"roles"`
	}

	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	return func(c echo.Context) error {
		var req request

		if err := c.Bind(&req); err != nil {
			return responses.Internal(c, err)
		}

		tokens, err := as.Register(c.Request().Context(), &dto.CreateUser{
			LastName:   req.LastName,
			FirstName:  req.FirstName,
			MiddleName: req.MiddleName,
			Email:      req.Email,
			Password:   req.Password,
			Roles:      req.Roles,
		})
		if err != nil {
			if errors.Is(err, authservice.ErrEmailTaken) {
				return responses.BadRequest(c, err)
			}
			return responses.Internal(c, err)
		}

		return c.JSON(200, &response{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
