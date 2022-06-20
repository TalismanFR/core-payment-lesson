package tests

import (
	"bytes"
	"context"
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	tc "github.com/testcontainers/testcontainers-go"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestIntergration(t *testing.T) {

	if err := godotenv.Load("../.env"); err != nil {
		panic("cannot load .env file")
	}

	apiHost, ok := os.LookupEnv("API_HOST")
	if !ok {
		panic("no API_HOST env")
	}
	shopId, ok := os.LookupEnv("SHOP_ID")
	if !ok {
		panic("no SHOP_ID env")
	}
	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		panic("no SECRET env")
	}

	t.Logf("API_HOST: %s\nSHOP_ID: %s\nSECRET: %s\n", apiHost, shopId, secret)

	dc := tc.NewLocalDockerCompose([]string{"./docker-compose.yaml"}, "my-id")

	dc.Executable = "/usr/local/bin/docker-compose"

	execError := dc.
		WithCommand([]string{"up", "-d"}).
		Invoke()
	if execError.Error != nil {
		t.Errorf("Failed when running: %v", execError.Command)
	}

	defer dc.Down()

	conf := config.Config{}

	// http://127.0.0.1:8200 terminals
	conf.Vault.Address = "http://127.0.0.1:8200"
	conf.Vault.MountPath = "terminals"

	// host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable
	conf.Payment.Host = "localhost"
	conf.Payment.User = "payservice"
	conf.Payment.Password = "payservice"
	conf.Payment.DBName = "payservice-db"
	conf.Payment.Port = "5432"
	conf.Payment.SslMode = "disable"

	// host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable
	conf.Terminal.Host = "localhost"
	conf.Terminal.User = "payservice"
	conf.Terminal.Password = "payservice"
	conf.Terminal.DBName = "payservice-db"
	conf.Terminal.Port = "5432"
	conf.Terminal.SslMode = "disable"

	time.Sleep(5 * time.Second)

	err := config.BuildDI(conf)
	if err != nil {
		t.Fatal(err)
	}

	insertStmt := fmt.Sprintf("INSERT INTO terminals(uuid, alias,url) VALUES ('8242df35-e182-4448-a99d-fd6b86dd7312','bepaid','%s')", apiHost)
	cmd := exec.Command("docker", []string{"exec", "my-id-postgres-1", "psql", "-U", "payservice", "-d", "payservice-db", "-c", insertStmt}...)

	var out, e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	t.Log("inserting terminal into repo")
	err = cmd.Run()
	t.Logf("insertion result: %q\n", out.String())
	t.Logf("insertion errors: %q\n", e.String())

	if err != nil {
		t.Fatal(err)
	}

	var ss application.SecretsRepository
	err = container.Resolve(&ss)
	if err != nil {
		t.Fatal("Resolve SecretsRepository", err)
	}

	err = ss.Put(context.Background(), uuid.MustParse("8242df35-e182-4448-a99d-fd6b86dd7312"), map[string]interface{}{"shop_id": shopId, "secret": secret})
	if err != nil {
		t.Fatalf("put error: %v\n", err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		t.Fatal(err)
	}

	cc := vo.NewCreditCard("4200000000000000", "123", "tim", vo.January, "2024")

	requestDto := *dto.NewChargeRequest(1000, "RUB", "8242df35-e182-4448-a99d-fd6b86dd7312", "invoiceID1", "description", *cc)

	t.Log("sending charge request")

	result, err := service.Charge(requestDto)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
