package charge

import (
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"diLesson/pkg/vault"
	"github.com/golobby/container/v3"
	tc "github.com/testcontainers/testcontainers-go"
	"gotest.tools/v3/env"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCharge_Integration(t *testing.T) {

	// Start postgres container
	dc := tc.NewLocalDockerCompose([]string{"docker-compose.yaml"}, "my-id")

	dc.Executable = "/usr/local/bin/docker-compose"

	execError := dc.WithCommand([]string{"up", "-d"}).Invoke()
	if execError.Error != nil {
		t.Errorf("Failed when running: %v", execError.Command)
	}

	defer dc.Down()

	// Create vault client
	v, err := vault.NewVault("http://127.0.0.1:8200", "terminals")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize and unseal vault
	// Get VAULT_ADDR and VAULT_TOKEN env
	envs, err := v.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	envs["POSTGRES_HOST"] = "localhost"
	envs["POSTGRES_PORT"] = "5432"
	envs["POSTGRES_USER"] = "payservice"
	envs["POSTGRES_PASSWORD"] = "payservice"

	// Setting envs
	unpatch := env.PatchAll(t, envs)
	defer unpatch()

	// Populating vault with values
	f, err := os.Open("terminals.json")
	if err != nil {
		t.Fatal(err)
	}

	uuids, err := v.Populate(f)
	if err != nil {
		t.Fatal(err)
	}

	// Uuids order is similar to order in terminals.json
	// Extract first terminal uuid (bepaid for now)
	terminalId := uuids[0]

	// Build config
	p, err := filepath.Abs("main.yaml")
	if err != nil {
		t.Fatal(err)
	}

	conf := config.Parse(p)

	time.Sleep(5 * time.Second)

	err = config.BuildDI(conf)
	if err != nil {
		t.Fatal(err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		t.Fatal(err)
	}

	cc := vo.NewCreditCard("4200000000000000", "123", "tim", vo.January, "2024")

	requestDto := *dto.NewChargeRequest(1000, "RUB", terminalId, "invoiceID1", "description", *cc)

	t.Log("sending charge request")

	result, err := service.Charge(requestDto)

	t.Log(result)

	if err == nil {
		t.Fatal("error shouldn't be nil")
	}

	if err.Error() != "Shop not found" {
		t.Fatalf("unexpected error message: AR: %v, ER: %s", err, "Shop not found")
	}
}
