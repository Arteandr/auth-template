package authservice

import (
	"context"
	"log/slog"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/lib/jwt"
	"mzhn/auth/internal/lib/logger/sl"
)

func (a *AuthService) Refresh(ctx context.Context, req *dto.Refresh) (*dto.Tokens, error) {
	log := a.logger.With("method", "AuthService.Refresh")

	log.Debug("refreshing", slog.Any("req", req))

	claims, err := jwt.Verify(req.RefreshToken, a.cfg.Jwt.RefreshSecret)
	if err != nil {
		log.Error("refresh token not valid", sl.Err(err))
		if err := a.sessionStorage.Delete(ctx, req.RefreshToken); err != nil {
			log.Error("delete session error", sl.Err(err))
			return nil, err
		}
		return nil, err
	}

	if err := a.sessionStorage.Check(ctx, claims.Id, req.RefreshToken); err != nil {
		log.Error("session not found", sl.Err(err))
		return nil, err
	}

	tokens, err := a.generateJwtPair(claims)
	if err != nil {
		log.Error("generate jwt pair error", sl.Err(err))
		return nil, err
	}

	if err := a.sessionStorage.Delete(ctx, claims.Id); err != nil {
		log.Error("delete session error", sl.Err(err))
		return nil, err
	}

	if err := a.sessionStorage.Save(ctx, claims.Id, tokens.RefreshToken); err != nil {
		log.Error("save session error", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
