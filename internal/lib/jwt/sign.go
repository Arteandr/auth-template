package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"mzhn/auth/internal/entity"
)

func Sign(user *entity.UserClaims, ttl time.Duration, secret []byte) (string, error) {
	payload := claims{
		UserClaims: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}
