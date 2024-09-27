package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"mzhn/auth/internal/entity"
)

type claims struct {
	entity.UserClaims
	jwt.RegisteredClaims
}
