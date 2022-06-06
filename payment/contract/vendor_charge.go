package contract

import (
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
)

type VendorCharge interface {
	Charge(pay *domain.Pay) (*dto.VendorChargeResult, error)
}
