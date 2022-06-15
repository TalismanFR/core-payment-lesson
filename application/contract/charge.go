package contract

import "payservice-core/application/contract/dto"

type Charge interface {
    Charge(request dto.ChargeRequest) (*dto.ChargeResult, error)
}
