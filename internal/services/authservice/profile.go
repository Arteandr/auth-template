package authservice

import (
	"context"
	"log/slog"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
)

func (a *AuthService) Profile(ctx context.Context, userId string) (*entity.User, []entity.Role, error) {

	log := a.logger.With(slog.String("method", "Profile"), slog.String("userId", userId))

	user, err := a.userStorage.Find(ctx, userId)
	if err != nil {
		log.Warn("cannot find user")
		return nil, nil, err
	}
	log.Debug("user found", slog.Any("user", user))

	roles, err := a.roleStorage.ListUser(ctx, userId)
	if err != nil {
		log.Warn("cannot list user's roles", sl.Err(err))
		return nil, nil, err
	}

	return user, roles, nil
}
