package contract

import (
    "payservice-core/application/domain"
    "payservice-core/payment/contract/dto"
)

type VendorCharge interface {
    Charge(pay *domain.Pay) (*dto.VendorChargeResult, error)
}
