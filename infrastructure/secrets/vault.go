package secrets

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"io/ioutil"
	"os"
	"strings"
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

func (v Vault) Validate() error {

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

		err = saveToFile("vault_secrets.txt", t)
		if err != nil {
			return err
		}

	} else {
		err := readFromFile("vault_secrets.txt", t)
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

func (v Vault) Get(ctx context.Context, terminalUuid uuid.UUID) (map[string]interface{}, error) {
	s, err := v.c.KVv2(v.mountPath).Get(ctx, terminalUuid.String())

	if err != nil {
		return nil, fmt.Errorf("unable to read secret data: %w", err)
	}

	return s.Data, nil
}

func (v Vault) Put(ctx context.Context, terminalUuid uuid.UUID, credentials map[string]interface{}) error {
	_, err := v.c.KVv2(v.mountPath).Put(ctx, terminalUuid.String(), credentials)

	if err != nil {
		return fmt.Errorf("unable to write secret: %w", err)
	}

	return nil
}

func saveToFile(fileName string, tokens *tokens) error {

	v := strings.Join(tokens.shared, " ")
	v = tokens.root + "\n" + v + "\n"

	return ioutil.WriteFile(fileName, []byte(v), 0644)
}

func readFromFile(fileName string, t *tokens) error {
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
