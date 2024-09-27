package config

import (
	"fmt"
	"log/slog"
	"os"

	"mzhn/auth/internal/lib/logger/prettyslog"
	"mzhn/auth/internal/lib/logger/sl"

	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Name    string `env:"APP_NAME" env-required:"true"`
	Version string `env:"APP_VERSION" env-required:"true"`
	Host    string `env:"APP_HOST" env-required:"true"`
	Port    int    `env:"APP_PORT" env-required:"true"`
}

type Pg struct {
	Host string `env:"PG_HOST" env-required:"true"`
	Port int    `env:"PG_PORT" env-required:"true"`
	User string `env:"PG_USER" env-required:"true"`
	Pass string `env:"PG_PASS" env-required:"true"`
	Name string `env:"PG_NAME" env-required:"true"`
}

type Redis struct {
	Host string `env:"REDIS_HOST" env-required:"true"`
	Port int    `env:"REDIS_PORT" env-required:"true"`
	Pass string `env:"REDIS_PASS"`
}

type Jwt struct {
	AccessSecret  string `env:"JWT_ACCESS_SECRET" env-required:"true"`
	AccessTTL     int    `env:"JWT_ACCESS_TTL" env-required:"true"`
	RefreshSecret string `env:"JWT_REFRESH_SECRET" env-required:"true"`
	RefreshTTL    int    `env:"JWT_REFRESH_TTL" env-required:"true"`
}

type Bcrypt struct {
	Cost int `env:"BCRYPT_COST" env-required:"true"`
}

type Config struct {
	Env    string `env:"ENV" env-default:"local"`
	App    App
	Pg     Pg
	Jwt    Jwt
	Bcrypt Bcrypt
	Redis  Redis
}

func New() *Config {
	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		slog.Error("error when reading env", sl.Err(err))
		header := fmt.Sprintf("%s - %s", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION"))

		usage := cleanenv.FUsage(os.Stdout, config, &header)
		usage()

		os.Exit(-1)
	}

	setupLogger(config)

	slog.Info("config", slog.Any("c", config))
	return config
}

func setupLogger(cfg *Config) {
	var log *slog.Logger

	switch cfg.Env {
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(prettyslog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	slog.SetDefault(log)
}
