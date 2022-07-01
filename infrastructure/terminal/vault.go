package terminal

import (
	"context"
	"diLesson/application/domain/terminal"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentation        = "infrastructure.terminal.vault"
	instrumentationVersion = "v0.0.1"
)

var (
	tracer = otel.Tracer(
		instrumentation,
		trace.WithSchemaURL(semconv.SchemaURL),
		trace.WithInstrumentationVersion(instrumentationVersion),
	)
)

type Vault struct {
	mountPath string
	c         *vault.Client
}

func NewVault(mountPath string) (*Vault, error) {
	config := vault.DefaultConfig()

	c, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vault{mountPath: mountPath, c: c}, nil
}

func (v Vault) FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*terminal.Terminal, error) {

	ctx, span := tracer.Start(ctx, "FindByUuid")
	defer span.End()

	s, err := v.c.KVv2(v.mountPath).Get(ctx, terminalUuid.String())

	if err != nil {
		return nil, fmt.Errorf("unable to read secret data: %w", err)
	}

	a, ok := s.Data["alias"]
	if !ok {
		return nil, fmt.Errorf("no alias key")
	}

	alias, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("alias has wrong type, expected: %T, got: %T", alias, a)
	}

	additionalData := map[string]interface{}{}

	for key, value := range s.Data {
		if key != "alias" {
			additionalData[key] = value
		}
	}

	return terminal.NewTerminal(terminalUuid, alias, additionalData), nil
}
