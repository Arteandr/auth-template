package pg

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
	"mzhn/auth/internal/services/authservice"
	"mzhn/auth/internal/storage"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

var _ authservice.UserStorage = (*UsersStorage)(nil)

type UsersStorage struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (s *UsersStorage) Find(ctx context.Context, slug string) (*entity.User, error) {
	log := s.logger.With(slog.String("user_id", slug)).With(slog.String("method", "Find"))

	builder := squirrel.Select().
		Columns("*").
		From(usersTable).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := uuid.Parse(slug); err != nil {
		slog.Debug("uuid parse error", sl.Err(err))
		builder = builder.Where(squirrel.Eq{"email": slug})
	} else {
		builder = builder.Where(squirrel.Eq{"id": slug})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("error building query", sl.Err(err))
		return nil, err
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	user := new(entity.User)
	err = s.db.GetContext(ctx, user, query, args...)
	if err != nil {
		log.Error("error to find user", sl.Err(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authservice.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UsersStorage) Save(ctx context.Context, user *dto.CreateUser) (*entity.User, error) {
	log := s.logger.With(slog.Any("user", user), slog.String("method", "Save"))

	builder := squirrel.
		Insert(usersTable).
		Columns("email", "hashed_password").
		Values(user.Email, user.Password).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)

	if user.FirstName != nil {
		builder.Columns("first_name").Values(user.FirstName)
	}

	if user.LastName != nil {
		builder.Columns("last_name").Values(user.LastName)
	}

	if user.MiddleName != nil {
		builder.Columns("middle_name").Values(user.MiddleName)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("error building query", sl.Err(err))
		return nil, err
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	newUser := new(entity.User)
	if err = s.db.GetContext(ctx, newUser, query, args...); err != nil {
		if e, ok := err.(pgx.PgError); ok {
			log.Debug("pg error", sl.PgError(e))
			if e.Code == "23505" {
				return nil, storage.ErrUserAlreadyExists
			}
		}
		log.Error("error saving user", sl.Err(err))
		return nil, err
	}

	return newUser, nil
}

func NewUserStorage(db *sqlx.DB) *UsersStorage {
	return &UsersStorage{
		db:     db,
		logger: slog.Default().With(slog.String("struct", "UserStorage")),
	}
}
