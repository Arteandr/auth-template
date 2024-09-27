package authservice

import (
	"context"
	"log/slog"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/jwt"
	"mzhn/auth/internal/lib/logger/sl"
)

func (a *AuthService) Authenticate(ctx context.Context, req *dto.Authenticate) (*entity.User, error) {

	log := a.logger.With("method", "Authenticate")
	log.Debug("authenticating", slog.Any("req", req))

	claims, err := jwt.Verify(req.AccessToken, a.cfg.Jwt.AccessSecret)
	if err != nil {
		log.Warn("invalid token", sl.Err(err))
		return nil, ErrTokenInvalid
	}

	log.Debug("claims", slog.Any("claims", claims))

	user, err := a.userStorage.Find(ctx, claims.Id)
	if err != nil {
		log.Error("user not found", sl.Err(err))
		return nil, ErrUserNotFound
	}

	ok, err := a.roleStorage.Check(ctx, &dto.CheckRoles{
		UserId: user.Id,
		Roles:  req.Roles,
	})
	if err != nil {
		log.Error("check roles error", sl.Err(err))
		return nil, err
	}

	log.Debug("authenticate result", slog.Any("ok", ok), slog.Any("user", user))

	if !ok {
		return nil, ErrInsufficientPermission
	}

	return user, nil
}
