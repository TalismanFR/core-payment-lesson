package terminal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
)

type BepaidShopCredentials struct {
	ShopId string
	Secret string
}

type TerminalSecrets interface {
	GetCredentials(ctx context.Context, uuid uuid.UUID) (*BepaidShopCredentials, error)
}

const (
	key1 = "shop_id"
	key2 = "secret"
)

type VaultTerminalSecrets struct {
	address   string //"http://127.0.0.1:8300"
	token     string //"myroot"
	mountPath string // "terminals"
}

func NewVaultTerminalSecrets(address string, token string, mountPath string) *VaultTerminalSecrets {
	return &VaultTerminalSecrets{address: address, token: token, mountPath: mountPath}
}

func (v VaultTerminalSecrets) GetCredentials(ctx context.Context, uuid uuid.UUID) (*BepaidShopCredentials, error) {
	config := vault.DefaultConfig()

	config.Address = v.address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	client.SetToken(v.token)

	kvSecret, err := client.KVv2(v.mountPath).Get(ctx, uuid.String())
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	shop_id, ok := kvSecret.Data[key1].(string)
	if !ok {
		return nil, fmt.Errorf("value type assertion failed: %T %#v", kvSecret.Data[key1], kvSecret.Data[key1])
	}

	secret, ok := kvSecret.Data[key2].(string)
	if !ok {
		return nil, fmt.Errorf("value type assertion failed: %T %#v", kvSecret.Data[key2], kvSecret.Data[key2])
	}

	return &BepaidShopCredentials{shop_id, secret}, nil
}
