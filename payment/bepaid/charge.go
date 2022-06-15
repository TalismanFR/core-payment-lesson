package bepaid

import (
	"fmt"
	"payservice-core/application/domain"
	"payservice-core/payment/contract/dto"
)

type Charge struct {
}

func (c Charge) Charge(pay *domain.Pay) (*dto.VendorChargeResult, error) {

    fmt.Println("charge service bepaid")
    return dto.NewVendorChargeResult("bepaid"), nil
}
