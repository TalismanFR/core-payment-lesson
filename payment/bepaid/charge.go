package bepaid

import (
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"fmt"
)

type Charge struct {
}

func (c Charge) Charge(pay *domain.Pay) (*dto.VendorChargeResult, error) {

	fmt.Println("charge service bepaid")
	return dto.NewVendorChargeResult("bepaid"), nil
}
