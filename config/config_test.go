package config

import (
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/env"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	unpatch := env.PatchAll(t, map[string]string{
		"POSTGRES_HOST":     "postgres",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "payservice",
		"POSTGRES_PASSWORD": "payservice",
		"VAULT_ADDR":        "example.com",
		"VAULT_TOKEN":       "12345",
	})

	defer unpatch()

	input := `http:
  port: 8000
  readTimeout: 10s
  writeTimeout: 5s

cache:
  ttl: 60s

auth:
  accessTokenTTL: 2h
  refreshTokenTTL: 720h #30 days
  verificationCodeLength: 8

limiter:
  rps: 10
  burst: 20
  ttl: 5m

postgres:
  dbName: payservice-db
  sslMode: disable

vault:
  mountPath: terminals
`
	p, _ := filepath.Abs("cfg1.yaml")
	f, err := os.CreateTemp(filepath.Dir(p), "cfg1.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fileName := f.Name()

	defer os.Remove(fileName) // clean up

	er := Config{
		Http: HttpConfig{
			Port:         "8000",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		Cache: CacheConfig{
			Ttl: 60 * time.Second,
		},
		Auth: AuthConfig{
			AccessTokenTTL:         2 * time.Hour,
			RefreshTokenTTL:        720 * time.Hour,
			VerificationCodeLength: 8,
		},
		Limiter: LimiterConfig{
			Rps:   10,
			Burst: 20,
			Ttl:   5 * time.Minute,
		},

		Postgres: PostgresConfig{
			Host:     "postgres",
			Port:     "5432",
			User:     "payservice",
			Password: "payservice",
			DBName:   "payservice-db",
			SslMode:  "disable",
		},

		Vault: VaultConfig{
			MountPath: "terminals",
		},
	}

	if _, err := f.Write([]byte(input)); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := Parse(fileName)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(er, *cfg) {
		t.Log(cmp.Diff(er, *cfg))
		t.Fatal()
	}
}
