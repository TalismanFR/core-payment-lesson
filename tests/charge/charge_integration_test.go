//go:build integration
// +build integration

package charge

import (
	"bufio"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"diLesson/pkg/vault"
	"github.com/golobby/container/v3"
	tc "github.com/testcontainers/testcontainers-go"
	"io"
	"os"
	"strings"
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

	v, err := vault.NewVault("http://127.0.0.1:8200", "terminals")
	if err != nil {
		t.Fatal(err)
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		if err = v.Initialize(pw); err != nil {
			t.Fatal(err)
		}
	}()

	sc := bufio.NewScanner(pr)
	for sc.Scan() {
		vs := strings.Split(sc.Text(), "=")
		err := os.Setenv(vs[0], vs[1])
		if err != nil {
			t.Fatal(err)
		}
	}

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
	conf := config.Config{}

	conf.Vault.MountPath = "terminals"

	// host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable
	conf.Payment.Host = "localhost"
	conf.Payment.User = "payservice"
	conf.Payment.Password = "payservice"
	conf.Payment.DBName = "payservice-db"
	conf.Payment.Port = "5432"
	conf.Payment.SslMode = "disable"

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
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
