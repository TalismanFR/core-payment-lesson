package contract

import (
	"diLesson/application/domain"
	"diLesson/payment/contract/dto"
)

type VendorRefund interface {
	Refund(pay *domain.Pay) (*dto.VendorRefundResult, error)
}
