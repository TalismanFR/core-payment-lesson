package config

import (
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/env"
	"os"
	"path"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestInit(t *testing.T) {
	tests := []struct {
		name      string
		inputYAML string
		inputEnv  map[string]string
		er        Config
	}{
		{
			name: "config file and environmental vars",
			inputYAML: `jwt:
  accessTokenTTL: 9s
  refreshTokenTTL: 3m

grpc:
  port: 50051

kafka:
  host: kafka
  port: 9003
  group: auth
  topic: newclient
  partition: 3

postgres:
  host: postgres
  port: 5555

`,
			inputEnv: map[string]string{
				"AUTH_POSTGRES_USER":     "authservice",
				"AUTH_POSTGRES_PASSWORD": "authservice",
				"AUTH_JWT_ACCESS_KEY":    "access_key",
				"AUTH_JWT_REFRESH_KEY":   "refresh_key",
			},
			er: Config{
				JWT: JWTConfig{
					AccessSigningKey:  "access_key",
					RefreshSigningKey: "refresh_key",
					AccessTokenTTL:    9 * time.Second,
					RefreshTokenTTL:   3 * time.Minute,
				},
				Grpc: GrpcConfig{
					Port: "50051",
				},
				Kafka: KafkaConfig{
					Host:      "kafka",
					Port:      "9003",
					Group:     "auth",
					Topic:     "newclient",
					Partition: "3",
				},
				Postgres: PostgresConfig{
					Host:     "postgres",
					Port:     "5555",
					User:     "authservice",
					Password: "authservice",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			unpatch := env.PatchAll(t, tc.inputEnv)
			defer unpatch()

			f, err := os.CreateTemp("", "config")
			defer os.Remove(f.Name())

			if _, err = f.WriteString(tc.inputYAML); err != nil {
				t.Fatalf("couldn't write to temp file %s: %s", f.Name(), err)
			}

			f.Close()

			cfg := Init(path.Join(f.Name()))

			if !cmp.Equal(tc.er, *cfg) {
				t.Fatal("configs aren't equal\n", cmp.Diff(tc.er, *cfg))
			}
		})
	}
}
