package pg

import (
	"context"
	"log/slog"

	"mzhn/auth/internal/dto"
	"mzhn/auth/internal/entity"
	"mzhn/auth/internal/lib/logger/sl"
	"mzhn/auth/internal/services/authservice"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

var _ authservice.RoleStorage = (*RoleStorage)(nil)

type RoleStorage struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewRoleStorage(db *sqlx.DB) *RoleStorage {
	return &RoleStorage{
		db:     db,
		logger: slog.Default().With(slog.String("struct", "RoleStorage")),
	}
}

func (r *RoleStorage) Add(ctx context.Context, dto *dto.AddRoles) (err error) {

	log := r.logger.With(slog.String("method", "Add"))
	log.Debug("dto", slog.Any("dto", dto))

	tx, err := r.db.Begin()
	if err != nil {
		log.Error("cannot begin transaction", sl.Err(err))
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return nil
		}

		if err := tx.Commit(); err != nil {
			log.Error("cannot commit transaction", sl.Err(err))
			return err
		}

		return nil
	}()

	for _, role := range dto.Roles {

		if !role.Valid() {
			log.Warn("invalid role", slog.String("role", role.String()))
			continue
		}

		query, args, err := squirrel.
			Insert(roleTable).
			Columns("user_id", "role").
			Values(dto.UserId, role).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			log.Error("cannot build query", sl.Err(err))
			return err
		}

		qlog := log.With(slog.String("query", query), slog.Any("args", args))

		qlog.Debug("executing query")

		if _, err := tx.Exec(query, args...); err != nil {
			qlog.Error("cannot execute query", sl.Err(err))
			return err
		}
	}

	return nil
}

func (r *RoleStorage) ListUser(ctx context.Context, userId string) ([]entity.Role, error) {
	log := r.logger.With(slog.String("method", "ListUser"))

	log.Debug("listing user's roles", slog.String("userId", userId))

	query, args, err := squirrel.
		Select("role").
		From(roleTable).
		Where(squirrel.Eq{"user_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return nil, err
	}

	qlog := log.With(slog.String("query", query), slog.Any("args", args))
	qlog.Debug("executing query")

	roles := make([]entity.Role, 0, 3)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		qlog.Error("cannot execute query", sl.Err(err))
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var role string

		if err := rows.Scan(&role); err != nil {
			qlog.Error("cannot scan row", sl.Err(err))
			return nil, err
		}

		roles = append(roles, entity.Role(role))
	}

	return roles, nil
}

func (r *RoleStorage) Check(ctx context.Context, dto *dto.CheckRoles) (bool, error) {
	log := r.logger.With(slog.String("method", "Check"))

	if len(dto.Roles) == 0 {
		return true, nil
	}

	log.Debug("dto", slog.Any("dto", dto))

	query, args, err := squirrel.
		Select("role").
		From(roleTable).
		Where(squirrel.Eq{"user_id": dto.UserId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return false, err
	}

	qlog := log.With(slog.String("query", query), slog.Any("args", args))
	qlog.Debug("executing query")

	rows, err := r.db.Query(query, args...)
	if err != nil {
		qlog.Error("cannot execute query", sl.Err(err))
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			qlog.Error("cannot scan row", sl.Err(err))
			return false, err
		}

		if lo.Contains(dto.Roles, entity.Role(role)) {
			return true, nil
		}
	}

	return false, nil
}

func (r *RoleStorage) Remove(ctx context.Context, dto *dto.RemoveRoles) (err error) {
	log := r.logger.With(slog.String("method", "Remove"))

	log.Debug("dto", slog.Any("dto", dto))

	tx, err := r.db.Begin()
	if err != nil {
		log.Error("cannot begin transaction", sl.Err(err))
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return nil
		}

		if err := tx.Commit(); err != nil {
			log.Error("cannot commit transaction", sl.Err(err))
			return err
		}
		return nil
	}()

	for _, role := range dto.Roles {
		query, args, err := squirrel.
			Delete(roleTable).
			Where(squirrel.Eq{"user_id": dto.UserId, "role": role}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			log.Error("cannot build query", sl.Err(err))
			return err
		}

		qlog := log.With(slog.String("query", query), slog.Any("args", args))

		qlog.Debug("executing query")

		if _, err := tx.Exec(query, args...); err != nil {
			qlog.Error("cannot execute query", sl.Err(err))
			return err
		}

	}

	return nil
}
