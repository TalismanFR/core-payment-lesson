package tinkoff

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
	"fmt"
)

type Charge struct {
}

func (c Charge) Charge(ctx context.Context, pay *domain.Pay) (*dto.VendorChargeResult, error) {

	fmt.Println("UNIMPLEMENTED charge service tinkoff")
	return &dto.VendorChargeResult{}, nil
}
