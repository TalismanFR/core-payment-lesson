package main

import (
	"context"
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"diLesson/server"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"log"
)

func main() {

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	panic("not implemented")

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

	err := config.BuildDI(conf)

	if err != nil {
		log.Fatal(err)
	}

	var ss application.SecretsRepository
	err = container.Resolve(&ss)
	if err != nil {
		log.Fatal(err)
	}

	err = ss.Put(context.Background(), uuid.MustParse("8242df35-e182-4448-a99d-fd6b86dd7312"), map[string]interface{}{"shop_id": "", "secret": ""})
	if err != nil {
		log.Fatalf("put error: %v\n", err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		log.Fatal(err)
	}

	cc := vo.NewCreditCard("4200000000000000", "123", "tim", vo.January, "2024")

	requestDto := *dto.NewChargeRequest(1000, "RUB", "8242df35-e182-4448-a99d-fd6b86dd7312", "invoiceID1", "decr", *cc)

	result, err := service.Charge(requestDto)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
