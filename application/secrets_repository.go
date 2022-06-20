package application

import (
	"context"
	"github.com/google/uuid"
)

type SecretsRepository interface {
	Get(ctx context.Context, key uuid.UUID) (map[string]interface{}, error)
	Put(ctx context.Context, key uuid.UUID, values map[string]interface{}) error
}
