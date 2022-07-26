package cmd

import (
	"auth/internal/auth/simple"
	"auth/internal/cache/redis"
	"auth/internal/config"
	"auth/internal/repo/pg"
	"auth/internal/server/v1"
	"auth/pkg/logger"
	"auth/pkg/token/jwt"
)

func Run(cfg *config.Config) error {

	loggerLevel, err := logger.LevelFromString(cfg.Logger.Level)
	if err != nil {
		return err
	}

	l := logger.New(loggerLevel)

	tokenAuthority, err := jwt.NewJWTokenAuthority(&cfg.JWT)
	if err != nil {
		return err
	}

	repo, err := pg.NewRepoPG(&cfg.Postgres, l)
	if err != nil {
		return err
	}

	c, err := redis.New(&cfg.Redis, logger.New(loggerLevel))
	if err != nil {
		return err
	}

	a := simple.New(tokenAuthority, repo, c, l)

	return v1.New(a, logger.New(logger.Debug)).Run(&cfg.Grpc)
}
