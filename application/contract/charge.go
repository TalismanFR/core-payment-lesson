package contract

import (
	"diLesson/application/contract/dto"
)

type Charge interface {
	Charge(request dto.ChargeRequest) (*dto.ChargeResult, error)
}
