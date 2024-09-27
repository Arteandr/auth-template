package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
	"mzhn/auth/internal/storage"
)

func (a *AuthService) Register(ctx context.Context, req *dto.CreateUser) (tokens *dto.Tokens, err error) {

	log := a.logger.With("method", "AuthService.Register")

	log.Debug("registering", slog.Any("req", req))

	log.Debug("hashing password", slog.String("password", req.Password))
	req.Password, err = a.hash(req.Password)
	if err != nil {
		log.Error("hash password error", sl.Err(err))
		return nil, err
	}

	log.Debug("creating user")

	user, err := a.userStorage.Save(ctx, req)
	if err != nil {
		log.Error("create user error", sl.Err(err))
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return nil, ErrEmailTaken
		}

		return nil, err
	}

	if err := a.roleStorage.Add(ctx, &dto.AddRoles{
		UserId: user.Id,
		Roles:  req.Roles,
	}); err != nil {
		log.Error("add roles error", sl.Err(err))
		return nil, fmt.Errorf("cannot add roles %w", err)
	}

	log.Debug("generating jwt pair")
	tokens, err = a.generateJwtPair(&entity.UserClaims{
		Id:    user.Id,
		Email: user.Email,
	})
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
