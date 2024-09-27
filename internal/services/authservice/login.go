package authservice

import (
	"context"
	"log/slog"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
)

func (a *AuthService) Login(ctx context.Context, req *dto.Login) (*dto.Tokens, error) {

	log := a.logger.With("method", "Login")

	log.Debug("logging in", slog.Any("req", req))

	user, err := a.userStorage.Find(ctx, req.Email)
	if err != nil {
		log.Error("user not found", sl.Err(err))
		return nil, err
	}

	if err := a.comparePassword(user.HashedPassword, req.Password); err != nil {
		log.Error("password not match", sl.Err(err))
		return nil, err
	}

	tokens, err := a.generateJwtPair(&entity.UserClaims{Id: user.Id, Email: user.Email})
	if err != nil {
		log.Error("generate jwt pair error", sl.Err(err))
		return nil, err
	}

	if a.sessionStorage.Save(ctx, user.Id, tokens.RefreshToken); err != nil {
		log.Error("save session error", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
