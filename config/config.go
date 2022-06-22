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
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func Parse(filePath string) *Config {

	f, err := os.Open(filePath)
	if err != nil {
		processError(err)
	}

	defer f.Close()

	var cfg Config

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		processError(err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		processError(err)
	}

	return &cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

type Config struct {

	// Vault address and token are stored in env and processed in vault library

	Http struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	} `yaml:"http"`

	Cache struct {
		Ttl time.Duration
	} `yaml:"cache"`

	Auth struct {
		AccessTokenTTL         time.Duration
		RefreshTokenTTL        time.Duration
		VerificationCodeLength int
	} `yaml:"auth"`

	Limiter struct {
		Rps   int
		Burst int
		Ttl   time.Duration
	} `yaml:"limiter"`

	Vault struct {
		// Address and token are taken from env VAULT_ADDR and VAULT_TOKEN

		// e.g. "terminals"
		MountPath string
	} `yaml:"vault"`

	Payment struct {
		Host     string `envconfig:"POSTGRES_HOST"`
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"dbName"`
		Port     string `envconfig:"POSTGRES_PORT"`
		SslMode  string `yaml:"sslMode"`
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
