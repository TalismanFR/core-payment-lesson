package contract

import (
	"context"
	"diLesson/application/contract/dto"
)

type Charge interface {
	//TODO: add context
	Charge(ctx context.Context, request dto.ChargeRequest) (*dto.ChargeResult, error)
}
