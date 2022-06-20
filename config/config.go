package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/infrastructure/secrets"
	"diLesson/infrastructure/terminal"
	"diLesson/payment/bepaid"
	contractPayment "diLesson/payment/contract"
	"diLesson/payment/tinkoff"
	"fmt"
	"github.com/golobby/container/v3"
)

type Config struct {
	Vault struct {
		Address   string
		MountPath string
	}

	Payment struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     string
		SslMode  string
	}

	Terminal struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     string
		SslMode  string
	}
}

func BuildDI(conf Config) (err error) {

	err = container.Transient(func() (application.SecretsRepository, error) {

		v, e := secrets.NewVault(conf.Vault.Address, conf.Vault.MountPath)
		if e != nil {
			return nil, e
		}

		if e = v.Validate(); e != nil {
			return nil, e
		}

		return v, nil
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

		pr, err := repository.NewPayRepositoryPgsql(dsn)
		return pr, err
	})

	err = container.NamedTransient("bepaid_charge", func() (contractPayment.VendorCharge, error) {
		s := bepaid.NewCharge("bepaid")

		return s, nil
	})

	err = container.NamedTransient("tinkoff_charge", func() (contractPayment.VendorCharge, error) {
		s := tinkoff.Charge{}

		return &s, nil
	})

	err = container.Transient(func() (application.TerminalRepo, error) {

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			conf.Terminal.Host,
			conf.Terminal.User,
			conf.Terminal.Password,
			conf.Terminal.DBName,
			conf.Terminal.Port,
			conf.Terminal.SslMode,
		)
		return terminal.NewRepoPG(dsn)
	})

	err = container.Transient(func() (contract.Charge, error) {
		var ch service.Charge

		err := container.Fill(&ch)

		return &ch, err
	})

	return err
}
