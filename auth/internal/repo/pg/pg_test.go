package pg

import (
	"auth/internal/config"
	"auth/internal/model/principal"
	"auth/pkg/logger"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestRepoPG_CreatePrincipal(t *testing.T) {

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

	r, err := NewRepoPG(&config.PostgresConfig{
		Host:     ip,
		Port:     mappedPort.Port(),
		User:     "authservice",
		Password: "authservice",
		DBname:   "authservice-db",
		Sslmode:  "disable",
	}, logger.New(logger.Debug))

	if err != nil {
		t.Fatal(err)
	}
	p := principal.Principal{
		Email:    "ex@mail.ru",
		Password: "12345",
		Role:     "user1",
	}

	p.HashPassword()
	p.Sanitize()

	err = r.Migrate(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	err = r.CreatePrincipal(context.TODO(), p)
	t.Log(err)

	p.Role = "registered"
	p.HashedPassword = []byte("ksdfjperiogij9")

	err = r.UpdatePrincipal(context.TODO(), p)
	t.Log(err)
	//
	p2, err := r.GetPrincipal(context.TODO(), p.Email)
	t.Log(err)
	t.Log(p2)
	t.Log(string(p2.HashedPassword))
}
