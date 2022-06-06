package main

import (
	"diLesson/application/contract"
	"diLesson/application/contract/dto"
	"diLesson/config"
	"fmt"
	"github.com/golobby/container/v3"
)

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

	requestDto := dto.ChargeRequest{
		Amount:      1000,
		TerminalId:  "1234",
		InvoiceId:   "1_test",
		Description: "",
	}
	result, err := service.Charge(requestDto)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}
