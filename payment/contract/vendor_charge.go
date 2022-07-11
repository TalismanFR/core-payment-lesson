package contract

import (
	"context"
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
)

type VendorCharge interface {
	Charge(ctx context.Context, pay *domain.Pay) (*dto.VendorChargeResult, error)
}
