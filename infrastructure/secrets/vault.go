package secrets

import (
	"bufio"
	"context"
	"diLesson/application"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"io/ioutil"
	"os"
	"strings"
)

type VaultService struct {
	mountPath string
	c         *vault.Client
}

func NewVaultService(address string, mountPath string) (*VaultService, error) {
	config := vault.DefaultConfig()
	config.Address = address

	c, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &VaultService{mountPath: mountPath, c: c}, nil
}

type tokens struct {
	root   string
	shared []string
}

func (v VaultService) Validate() error {

	sys := v.c.Sys()

	s, err := sys.SealStatus()
	if err != nil {
		return err
	}

	isInit, isSealed := s.Initialized, s.Sealed

	t := &tokens{}

	if !isInit {
		initR := &vault.InitRequest{SecretShares: 1, SecretThreshold: 1}

		resp, err := sys.Init(initR)
		if err != nil {
			return err
		}

		t.shared = resp.KeysB64
		t.root = resp.RootToken

		err = SaveToFile("vault_secrets.txt", t)
		if err != nil {
			return err
		}

	} else {
		err := ReadFromFile("vault_secrets.txt", t)
		if err != nil {
			return err
		}
	}

	v.c.SetToken(t.root)

	if isSealed {
		_, err := sys.Unseal(t.shared[0])
		if err != nil {
			return err
		}
	}

	mounts, err := sys.ListMounts()
	if err != nil {
		return err
	}

	if _, ok := mounts[v.mountPath+"/"]; !ok {
		err = sys.Mount(v.mountPath, &vault.MountInput{Type: "kv", Options: map[string]string{"version": "2"}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (v VaultService) Get(ctx context.Context, terminalUuid uuid.UUID) (*application.BepaidShopCredentials, error) {
	s, err := v.c.KVv2(v.mountPath).Get(ctx, terminalUuid.String())
	if err != nil {
		return nil, fmt.Errorf("unable to read secret data: %w", err)
	}

	v1, ok := s.Data["shop_id"]
	if !ok {
		return nil, fmt.Errorf("enable to read secret shop_id")
	}

	shop_id, ok := v1.(string)
	if !ok {
		return nil, fmt.Errorf("shop_id type isn't a string. type: %T", shop_id)
	}

	v2, ok := s.Data["secret"]
	if !ok {
		return nil, fmt.Errorf("enable to read secret secret")
	}

	secret, ok := v2.(string)
	if !ok {
		return nil, fmt.Errorf("secret type isn't a string. type: %T", secret)
	}

	return &application.BepaidShopCredentials{ShopId: shop_id, Secret: secret}, nil
}

func (v VaultService) Put(ctx context.Context, terminalUuid uuid.UUID, credentials *application.BepaidShopCredentials) error {
	data := map[string]interface{}{"shop_id": credentials.ShopId, "secret": credentials.Secret}
	_, err := v.c.KVv2(v.mountPath).Put(ctx, terminalUuid.String(), data)

	if err != nil {
		return fmt.Errorf("unable to write secret: %w", err)
	}

	return nil
}

func SaveToFile(fileName string, tokens *tokens) error {

	v := strings.Join(tokens.shared, " ")
	v = tokens.root + "\n" + v + "\n"

	return ioutil.WriteFile(fileName, []byte(v), 0644)
}

func ReadFromFile(fileName string, t *tokens) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		return fmt.Errorf("file ends unexpectedly")
	}

	t.root = sc.Text()

	shared := make([]string, 0, 6)
	for sc.Scan() {
		shared = append(shared, sc.Text())
	}

	t.shared = shared

	return nil
}
