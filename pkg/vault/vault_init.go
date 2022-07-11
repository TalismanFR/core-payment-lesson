package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"io"
)

type Vault struct {
	mountPath string
	c         *vault.Client
}

func NewVault(address string, mountPath string) (*Vault, error) {
	config := vault.DefaultConfig()
	config.Address = address

	c, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vault{mountPath: mountPath, c: c}, nil
}

type tokens struct {
	root   string
	shared []string
}

func (v Vault) Initialize() (map[string]string, error) {

	envs := map[string]string{}

	sys := v.c.Sys()

	s, err := sys.SealStatus()
	if err != nil {
		return nil, err
	}

	isInit, isSealed := s.Initialized, s.Sealed

	t := &tokens{}

	if isInit {
		return nil, fmt.Errorf("vault shouldn't be initialized")
	}

	initR := &vault.InitRequest{SecretShares: 1, SecretThreshold: 1}

	resp, err := sys.Init(initR)
	if err != nil {
		return nil, err
	}

	t.shared = resp.KeysB64
	t.root = resp.RootToken

	v.c.SetToken(t.root)

	if isSealed {
		_, err := sys.Unseal(t.shared[0])
		if err != nil {
			return nil, err
		}
	}

	mounts, err := sys.ListMounts()
	if err != nil {
		return nil, err
	}

	if _, ok := mounts[v.mountPath+"/"]; !ok {
		err = sys.Mount(v.mountPath, &vault.MountInput{Type: "kv", Options: map[string]string{"version": "2"}})
		if err != nil {
			return nil, err
		}
	}

	envs["VAULT_TOKEN"] = t.root
	envs["VAULT_ADDR"] = v.c.Address()

	return envs, nil
}

// Populate decodes reader as json object and puts values in v.MountPath
func (v *Vault) Populate(in io.ReadCloser) (map[string]string, error) {
	defer in.Close()

	uuids := map[string]string{}

	m := map[string]interface{}{}

	if err := json.NewDecoder(in).Decode(&m); err != nil {
		return nil, err
	}

	ctx := context.Background()

	if len(m) == 0 {
		return nil, fmt.Errorf("no terminals in file")
	}

	for alias, data := range m {
		creds, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("wrong data type: expect %T, got: %T", creds, data)
		}

		creds["alias"] = alias
		uid := uuid.New().String()
		uuids[alias] = uid
		_, err := v.c.KVv2(v.mountPath).Put(ctx, uid, creds)
		if err != nil {
			return nil, fmt.Errorf("unable to write secret: %w", err)
		}
	}

	return uuids, nil
}
