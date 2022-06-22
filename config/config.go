package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/infrastructure/terminal"
	"diLesson/payment/bepaid"
	contractPayment "diLesson/payment/contract"
	"diLesson/payment/tinkoff"
	"fmt"
	"github.com/golobby/container/v3"
	"time"
)

type Config struct {

	// Vault address and token are stored in env and processed in vault library

	Http struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Cache struct {
		Ttl time.Duration
	}

	Auth struct {
		AccessTokenTTL         time.Duration
		RefreshTokenTTL        time.Duration
		VerificationCodeLength int
	}

	Limiter struct {
		Rps   int
		Burst int
		Ttl   time.Duration
	}

	Vault struct {
		// e.g. "terminals"
		MountPath string
	}

	Payment struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     string
		SslMode  string
	} `yaml:"postgres"`
}

func BuildDI(conf Config) (err error) {

	err = container.Transient(func() (application.TerminalRepo, error) {
		return terminal.NewVault(conf.Vault.MountPath)
	})

	err = container.Transient(func() (application.PayRepository, error) {

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			conf.Payment.Host,
			conf.Payment.User,
			conf.Payment.Password,
			conf.Payment.DBName,
			conf.Payment.Port,
			conf.Payment.SslMode,
		)

		return repository.NewPayRepositoryPgsql(dsn)
	})

	err = container.NamedTransient("bepaid_charge", func() (contractPayment.VendorCharge, error) {
		s := bepaid.NewCharge("bepaid")

		return s, nil
	})

	err = container.NamedTransient("tinkoff_charge", func() (contractPayment.VendorCharge, error) {
		s := tinkoff.Charge{}

		return &s, nil
	})

	err = container.Transient(func() (contract.Charge, error) {
		var ch service.Charge

		err := container.Fill(&ch)

		return &ch, err
	})

	return err
}
