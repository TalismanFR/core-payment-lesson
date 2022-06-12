package main

import (
	"context"
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/application/domain/vo"
	"diLesson/config"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"log"
)

func main() {

	err := config.BuildDI()

	if err != nil {
		log.Fatal(err)
	}

	var ss application.SecretsService
	err = container.Resolve(&ss)
	if err != nil {
		log.Fatal(err)
	}

	err = ss.Put(context.Background(), uuid.MustParse("8242df35-e182-4448-a99d-fd6b86dd7312"), &application.BepaidShopCredentials{ShopId: "19616", Secret: "877a3f17c9d196244dace3ab73b12cb89cafd77312c4301c2a71a02a653d1b9e"})
	if err != nil {
		log.Fatalf("put error: %v\n", err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		log.Fatal(err)
	}

	requestDto := *dto.NewChargeRequest(1000, "RUB", "8242df35-e182-4448-a99d-fd6b86dd7312", "invoiceID1", "decr", vo.CreditCard{})

	result, err := service.Charge(requestDto)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
