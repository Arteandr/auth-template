package authservice

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/auth/internal/lib/logger/sl"
)

func (a *AuthService) Logout(ctx context.Context, userId string) error {

	log := a.logger.With(slog.String("method", "Logout"))

	if _, err := a.userStorage.Find(ctx, userId); err != nil {
		log.Warn("user not found to logout", sl.Err(err))
		return fmt.Errorf("user not exists %w", err)
	}

	if err := a.sessionStorage.Delete(ctx, userId); err != nil {
		log.Warn("failed to delete session", sl.Err(err))
		return fmt.Errorf("failed to delete session %w", err)
	}

	return nil
}
