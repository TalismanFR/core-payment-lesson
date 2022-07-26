package v1

import (
	"auth/internal/api/v1"
	"auth/internal/auth"
	"auth/internal/auth/simple"
	"auth/internal/cache/memory"
	"auth/internal/config"
	"auth/internal/mocks"
	"auth/internal/repo/pg"
	"auth/pkg/logger"
	"auth/pkg/token/jwt"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

func TestService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewMockAuthenticationService(ctrl)

	tests := []struct {
		name        string
		input       v1.SignUpRequest
		authService func() auth.AuthenticationService
		isErr       bool
		err         error //TODO: domain level errors
	}{
		{
			name: "valid_1",
			input: v1.SignUpRequest{
				Email:    "",
				Password: "",
				Role:     "",
			},
		},
	}

	_ = tests

	ctx := context.TODO()
	mockAuth.EXPECT().Signup(ctx, auth.SignUpRequest{
		Email:    "",
		Password: "",
		Role:     "",
	}).Return(auth.JWTokens{
		Access:  "",
		Refresh: "",
	}, nil)

	res, err := New(mockAuth, logger.New(logger.Debug)).SignUp(ctx, &v1.SignUpRequest{
		Email:    "",
		Password: "",
		Role:     "",
	})

	if err != nil {
		t.Fatal(err)
	}

	_ = res

}

func TestServer_SignupWithTheSameCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewMockAuthenticationService(ctrl)

	ctx := context.TODO()
	email := "tim@ya.ru"
	password := "password123"
	role := "user"

	expectedTokens := auth.JWTokens{
		Access:  "AT1",
		Refresh: "RT1",
	}

	gomock.InOrder(
		mockAuth.EXPECT().Signup(ctx, auth.SignUpRequest{
			Email:    email,
			Password: password,
			Role:     role,
		}).Return(&expectedTokens, nil),
		mockAuth.EXPECT().Signup(ctx, auth.SignUpRequest{
			Email:    email,
			Password: password,
			Role:     role,
		}).Return(nil, errors.New("email taken")),
	)

	s := New(mockAuth, logger.New(logger.Debug))

	actualTokens, err := s.SignUp(ctx, &v1.SignUpRequest{
		Email:    email,
		Password: password,
		Role:     role,
	})

	require.Nil(t, err, "unexpected error: %v", err)
	require.Equal(t, expectedTokens, auth.JWTokens{
		Access:  actualTokens.GetAccess(),
		Refresh: actualTokens.GetRefresh(),
	})

	actualTokens, err = s.SignUp(ctx, &v1.SignUpRequest{
		Email:    email,
		Password: password,
		Role:     role,
	})

	require.NotNil(t, err, "second signup request with the same credentials should return non-nil error, got nil")
}

func TestIntegrationServer(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.3-alpine3.16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "authservice-db",
			"POSTGRES_USER":     "authservice",
			"POSTGRES_PASSWORD": "authservice",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatal(err)
	}

	ip, err := container.Host(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	mappedPort, err := container.MappedPort(context.TODO(), "5432")
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Config{
		JWT: config.JWTConfig{
			AccessSigningKey:  "123",
			RefreshSigningKey: "456",
			AccessTokenTTL:    5 * time.Second,
			RefreshTokenTTL:   10 * time.Second,
		},
		Grpc: config.GrpcConfig{
			Port: "50051",
		},
		Postgres: config.PostgresConfig{
			Host:     ip,
			Port:     mappedPort.Port(),
			User:     "authservice",
			Password: "authservice",
			DBname:   "authservice-db",
			Sslmode:  "disable",
		},
	}

	log.Println("start server")
	err = run(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("end server")

}

func run(cfg *config.Config) error {

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

	err = repo.Migrate(context.TODO())
	if err != nil {
		return err
	}

	// TODO: add to config
	cacheExpirationTime := 10 * time.Second
	c := memory.New(cacheExpirationTime)

	a := simple.New(tokenAuthority, repo, c, l)

	return New(a, logger.New(logger.Debug)).Run(&cfg.Grpc)
}
