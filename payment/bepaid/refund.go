package bepaid

import (
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"fmt"
)

type Refund struct {
}

func (c Charge) Refund(pay *domain.Pay) (*dto.VendorRefundResult, error) {

	fmt.Println("charge service bepaid")
	return dto.NewVendorRefundResult("bepaid"), nil
}
