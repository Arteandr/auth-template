//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"

	"mzhn/auth/internal/config"
	"mzhn/auth/internal/services/authservice"
	"mzhn/auth/internal/storage/pg"

	rd "mzhn/auth/internal/storage/redis"

	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		pg.NewUserStorage,
		pg.NewRoleStorage,
		rd.NewSessionsStorage,

		authservice.New,

		initPG,
		initRedis,
		config.New,

		wire.Bind(new(authservice.RoleStorage), new(*pg.RoleStorage)),
		wire.Bind(new(authservice.UserStorage), new(*pg.UsersStorage)),
		wire.Bind(new(authservice.SessionsStorage), new(*rd.SessionsStorage)),
	))
}

func initPG(cfg *config.Config) (*sqlx.DB, func(), error) {
	host := cfg.Pg.Host
	port := cfg.Pg.Port
	user := cfg.Pg.User
	pass := cfg.Pg.Pass
	name := cfg.Pg.Name

	cs := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`, user, pass, host, port, name)

	slog.Info("connecting to database", slog.String("conn", cs))

	db, err := sqlx.Connect("pgx", cs)
	if err != nil {
		return nil, nil, err
	}

	slog.Info("send ping to database")

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", slog.String("err", err.Error()), slog.String("conn", cs))
		return nil, func() { db.Close() }, err
	}

	slog.Info("connected to database", slog.String("conn", cs))

	return db, func() { db.Close() }, nil
}

func initRedis(cfg *config.Config) (*redis.Client, func(), error) {
	host := cfg.Redis.Host
	port := cfg.Redis.Port
	pass := cfg.Redis.Pass

	cs := fmt.Sprintf(`redis://%s:%s@%s:%d`, host, pass, host, port)

	slog.Info("connecting to redis", slog.String("conn", cs))

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		slog.Error("failed to connect to redis", slog.String("err", err.Error()), slog.String("conn", cs))
		return nil, func() { client.Close() }, err
	}

	slog.Info("connected to redis", slog.String("conn", cs))

	return client, func() {
		client.Close()
	}, nil
}
