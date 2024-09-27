package authservice

import (
	"context"
	"log/slog"

	"mzhn/auth/internal/config"
	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
)

type UserStorage interface {
	Find(ctx context.Context, slug string) (*entity.User, error)
	Save(ctx context.Context, user *dto.CreateUser) (*entity.User, error)
}

type SessionsStorage interface {
	Check(ctx context.Context, userId, token string) error
	Save(ctx context.Context, userId, token string) error
	Delete(ctx context.Context, userId string) error
}

type RoleStorage interface {
	Check(ctx context.Context, dto *dto.CheckRoles) (bool, error)
	ListUser(ctx context.Context, userId string) ([]entity.Role, error)
	Add(ctx context.Context, dto *dto.AddRoles) error
	Remove(ctx context.Context, dto *dto.RemoveRoles) error
}

type AuthService struct {
	userStorage    UserStorage
	roleStorage    RoleStorage
	sessionStorage SessionsStorage
	cfg            *config.Config
	logger         *slog.Logger
}

func New(userStorage UserStorage, roleStorage RoleStorage, sessionStorage SessionsStorage, cfg *config.Config) *AuthService {
	return &AuthService{
		cfg:            cfg,
		userStorage:    userStorage,
		roleStorage:    roleStorage,
		sessionStorage: sessionStorage,
		logger:         slog.Default().With(slog.String("struct", "AuthService")),
	}
}
