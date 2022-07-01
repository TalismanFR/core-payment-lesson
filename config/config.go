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

type (
	Config struct {
		Http     HttpConfig     `yaml:"http"`
		Grpc     GrpcConfig     `yaml:"grpc"`
		Cache    CacheConfig    `yaml:"cache"`
		Auth     AuthConfig     `yaml:"auth"`
		Limiter  LimiterConfig  `yaml:"limiter"`
		Postgres PostgresConfig `yaml:"postgres"`
		Vault    VaultConfig    `yaml:"vault"`
		Trace    TraceConfig    `yaml:"trace"`
	}

	HttpConfig struct {
		Port         string
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
	}

	GrpcConfig struct {
		Port string `yaml:"port"`
	}

	CacheConfig struct {
		Ttl time.Duration
	}

	AuthConfig struct {
		AccessTokenTTL         time.Duration `yaml:"accessTokenTTL"`
		RefreshTokenTTL        time.Duration `yaml:"refreshTokenTTL"`
		VerificationCodeLength int           `yaml:"verificationCodeLength"`
	}

	LimiterConfig struct {
		Rps   int
		Burst int
		Ttl   time.Duration
	}

	VaultConfig struct {
		MountPath string `yaml:"mountPath"`
	}

	PostgresConfig struct {
		Host     string `envconfig:"POSTGRES_HOST"`
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"dbName"`
		Port     string `envconfig:"POSTGRES_PORT"`
		SslMode  string `yaml:"sslMode"`
	}

	TraceConfig struct {
		Host     string `yaml:"host"`
		GrpcPort string `yaml:"grpcPort"`
		HttpPort string `yaml:"httpPort"`
	}
)

func Parse(filePath string) (*Config, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	cfg := &Config{}

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func BuildDI(conf *Config) (err error) {

	err = container.Transient(func() (application.TerminalRepo, error) {
		return terminal.NewVault(conf.Vault.MountPath)
	})

	err = container.Transient(func() (application.PayRepository, error) {

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			conf.Postgres.Host,
			conf.Postgres.User,
			conf.Postgres.Password,
			conf.Postgres.DBName,
			conf.Postgres.Port,
			conf.Postgres.SslMode,
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
