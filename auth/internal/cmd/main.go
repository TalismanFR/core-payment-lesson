package cmd

import (
	"auth/internal/auth/simple"
	"auth/internal/cache/memory"
	"auth/internal/config"
	"auth/internal/repo/pg"
	"auth/internal/server/v1"
	"auth/pkg/logger"
	"auth/pkg/token/jwt"
	"time"
)

func Run(cfg *config.Config) error {

	// TODO: add to config
	l := logger.New(logger.Debug)

	tokenAuthority, err := jwt.NewJWTokenAuthority(&cfg.JWT)
	if err != nil {
		return err
	}

	repo, err := pg.NewRepoPG(&cfg.Postgres, l)
	if err != nil {
		return err
	}

	// TODO: add to config
	cacheExpirationTime := 10 * time.Second
	c := memory.New(cacheExpirationTime)

	a := simple.New(tokenAuthority, repo, c, l)

	return v1.New(a, logger.New(logger.Debug)).Run(&cfg.Grpc)
}
