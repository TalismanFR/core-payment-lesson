package application

import (
	"context"
	"github.com/google/uuid"
)

type BepaidShopCredentials struct {
	ShopId string
	Secret string
}

type SecretsService interface {
	Get(ctx context.Context, terminalUuid uuid.UUID) (*BepaidShopCredentials, error)
	Put(ctx context.Context, terminalUuid uuid.UUID, credentials *BepaidShopCredentials) error
}
