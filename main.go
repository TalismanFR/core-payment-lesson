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
		TerminalId:  "terminalId1",
		InvoiceId:   "1_test",
		Description: "",
	}
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
