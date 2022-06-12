package main

import (
	"context"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"fmt"
	"github.com/golobby/container/v3"
	vault "github.com/hashicorp/vault/api"
	"log"
)

func main1() {
	config := vault.DefaultConfig()

	config.Address = "http://127.0.0.1:8300"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("myroot")

	//secretData := map[string]interface{}{
	//	"password": "Hashi123",
	//}
	//
	//// Write a secret
	//_, err = client.KVv2("secret").Put(context.Background(), "my-secret-password", secretData)
	//if err != nil {
	//	log.Fatalf("unable to write secret: %v", err)
	//}
	//
	//fmt.Println("Secret written successfully.")

	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("terminals").Get(context.Background(), "8242df35-e182-4448-a99d-fd6b86dd7312")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	value, ok := secret.Data["shop_id"].(string)
	if !ok {
		log.Fatalf("value type assertion failed: %T %#v", secret.Data["shop_id"], secret.Data["shop_id"])
	}

	if value != "shop_id" {
		log.Fatalf("unexpected password value %q retrieved from vault", value)
	}

	fmt.Println("Access granted!")
}

func main() {

	err := config.BuildDI()

	if err != nil {
		fmt.Println(err)
		return
	}

	var service contract.Charge

	err = container.Resolve(&service)

	if err != nil {
		fmt.Println(err)
		return
	}

	requestDto := *dto.NewChargeRequest(1000, "RUB", "8242df35-e182-4448-a99d-fd6b86dd7312", "invoiceID1", "decr", vo.CreditCard{})

	result, err := service.Charge(requestDto)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	//requestDto.Amount += 1
	//
	//err = service.Update(uuid.MustParse("acfad81d-03e0-4a9c-9d63-35e47b682a42"), requestDto)
	//if err != nil {
	//	fmt.Println(err)
	//}

}
