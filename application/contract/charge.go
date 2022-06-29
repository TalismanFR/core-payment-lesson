package contract

import (
	"diLesson/application/contract/dto"
)

type Charge interface {
	//TODO: add context
	Charge(request dto.ChargeRequest) (*dto.ChargeResult, error)
}
