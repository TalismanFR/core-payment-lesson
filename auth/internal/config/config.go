package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const envPrefix = "AUTH"

// TODO: add config validation

type (
	Config struct {
		JWT      JWTConfig `envconfig:"JWT"`
		Grpc     GrpcConfig
		Kafka    KafkaConfig
		Postgres PostgresConfig `envconfig:"POSTGRES"`
		Redis    RedisConfig    `envconfig:"REDIS"`
		Logger   Logger
	}

	JWTConfig struct {
		AccessSigningKey  string        `envconfig:"ACCESS_KEY"`
		RefreshSigningKey string        `envconfig:"REFRESH_KEY"`
		AccessTokenTTL    time.Duration `yaml:"accessTokenTTL"`
		RefreshTokenTTL   time.Duration `yaml:"refreshTokenTTL"`
	}

	GrpcConfig struct {
		Port string `yaml:"port"`
	}

	KafkaConfig struct {
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Group     string `yaml:"group"`
		Topic     string `yaml:"topic"`
		Partition string `yaml:"partition"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `envconfig:"user"`
		Password string `envconfig:"password"`
		DBname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	}

	RedisConfig struct {
		ExpireIn time.Duration `yaml:"expireIn"`
		Host     string        `yaml:"host"`
		Port     string        `yaml:"port"`
		Password string        `envconfig:"PASSWORD"`
		DB       int           `yaml:"db"`
	}

	Logger struct {
		Level string `yaml:"level"`
	}
)

func Init(file string) *Config {

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg := &Config{}

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		panic(err)
	}

	envconfig.MustProcess(envPrefix, cfg)

	return cfg
}
