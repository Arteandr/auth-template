package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"mzhn/auth/internal/config"
	"mzhn/auth/internal/lib/logger/sl"
	"mzhn/auth/internal/services/authservice"

	"github.com/redis/go-redis/v9"
)

var _ authservice.SessionsStorage = (*SessionsStorage)(nil)

type SessionsStorage struct {
	db     *redis.Client
	cfg    *config.Config
	logger *slog.Logger
}

func (s *SessionsStorage) Save(ctx context.Context, userId string, token string) error {
	log := s.logger.With(slog.String("user_id", userId))

	log.Debug("Saving session")

	if stat := s.db.Set(ctx, userId, token, time.Duration(s.cfg.Jwt.RefreshTTL)*time.Minute); stat.Err() != nil {
		log.Error("error saving session", sl.Err(stat.Err()))
		return fmt.Errorf("failed saving session %w", stat.Err())
	}

	return nil
}

func (s *SessionsStorage) Check(ctx context.Context, userId, refreshToken string) error {
	log := s.logger.With(slog.String("user_id", userId))

	log.Debug("Checking session")

	stat := s.db.Get(ctx, userId)
	if err := stat.Err(); err != nil {
		log.Error("error checking session", sl.Err(err))
		return fmt.Errorf("failed checking session %w", err)
	}

	if stat.Val() != refreshToken {
		log.Error("invalid session", slog.String("user_id", userId), slog.String("refresh_token", refreshToken), slog.String("session_token", stat.Val()))
		return fmt.Errorf("invalid session")
	}

	return nil
}

func (s *SessionsStorage) Delete(ctx context.Context, userId string) error {
	log := s.logger.With(slog.String("method", "SessionStorage.Delete"))

	log.Debug("Deleting session")

	stat := s.db.Del(ctx, userId)
	if err := stat.Err(); err != nil {
		log.Error("error deleting session", sl.Err(err))
		return fmt.Errorf("failed deleting session %w", err)
	}

	return nil
}

func NewSessionsStorage(db *redis.Client, cfg *config.Config) *SessionsStorage {
	return &SessionsStorage{
		db:     db,
		cfg:    cfg,
		logger: slog.Default().With(slog.String("struct", "UserStorage")),
	}
}
