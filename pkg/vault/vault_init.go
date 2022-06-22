package vault

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"io"
	"log"
	"os"
)

// Initialize and unseal vault, add kv2 engine, populate with values from file
func main() {

	var addr, mountPath, file, out string
	flag.StringVar(&addr, "address", "http://127.0.0.1:8200", "vault address")
	flag.StringVar(&mountPath, "mount", "terminals", "engine name")
	flag.StringVar(&file, "file", "terminals.json", "file with terminals credentials")
	flag.StringVar(&out, "out", ".env", "file to write VAULT_ADDRESS and VAULT_TOKEN")

	flag.Parse()

	fmt.Println("args: ", addr, mountPath, file)

	v, err := NewVault(addr, mountPath)
	if err != nil {
		log.Fatal("NewVault error: ", err)
	}

	f := os.Stdout
	if out != "" {
		f, err = os.Create(out)
		defer f.Close()

		if err != nil {
			log.Fatalf("open file %s error: %s\n", out, err)
		}
	}

	if err := v.Initialize(f); err != nil {
		log.Fatal("Vault.Initialize error: ", err)
	}

	f2, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file %s error: %v\n", file, err)
	}

	if _, err := v.Populate(f2); err != nil {
		log.Fatal("Vault.Populate error : ", err)
	}
}

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

func (v Vault) Initialize(out io.Writer) error {

	sys := v.c.Sys()

	s, err := sys.SealStatus()
	if err != nil {
		return err
	}

	isInit, isSealed := s.Initialized, s.Sealed

	t := &tokens{}

	if isInit {
		return fmt.Errorf("vault shouldn't be initialized")
	}

	initR := &vault.InitRequest{SecretShares: 1, SecretThreshold: 1}

	resp, err := sys.Init(initR)
	if err != nil {
		return err
	}

	t.shared = resp.KeysB64
	t.root = resp.RootToken

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

	_, err = fmt.Fprintf(out, "VAULT_TOKEN=%s\n", t.root)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "VAULT_ADDR=%s\n", v.c.Address())
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Populate(in io.ReadCloser) ([]string, error) {
	defer in.Close()

	var ids []string

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
		ids = append(ids, uid)
		_, err := v.c.KVv2(v.mountPath).Put(ctx, uid, creds)
		if err != nil {
			return nil, fmt.Errorf("unable to write secret: %w", err)
		}
	}

	return ids, nil
}
