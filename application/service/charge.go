package service

import (
	"context"
	"diLesson/application"
	"diLesson/application/contract/dto"
	"diLesson/application/domain"
	"diLesson/application/domain/vo"
	"diLesson/payment/contract"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"time"
)

//func init() {
//	if err := config.BuildDI(); err != nil {
//		panic(fmt.Sprintf("couldn't build dependencies for the application: %v", err))
//	}
//}

type Charge struct {
	payRepo      application.PayRepository `container:"type"`
	terminalRepo application.TerminalRepo  `container:"type"`
}

func (c Charge) Update(id uuid.UUID, request dto.ChargeRequest) error {
	return c.payRepo.Update(context.Background(), domain.NewPay(id, vo.Amount(request.Amount), vo.RUB, request.InvoiceId, domain.StatusCodeOK, vo.StatusNew, time.Now(), "transactionId"))
}

func (c Charge) Charge(request dto.ChargeRequest) (*dto.ChargeResult, error) {

	if request.Amount < 0 {
		return nil, fmt.Errorf("amount value less than zero")
	}

	// TODO: when to generate and when to acquire transactionId
	transactionId := ""

	pay := domain.NewPay(uuid.New(), vo.Amount(request.Amount), vo.RUB, request.InvoiceId, domain.StatusCodeOK, vo.StatusNew, time.Now(), transactionId)

	a, err := c.terminalRepo.GetAlias(request.TerminalId)
	if err != nil {
		return nil, err
	}

	var vendor contract.VendorCharge
	container.NamedResolve(&vendor, a)

	result, err := vendor.Charge(pay)
	if err != nil {
		return nil, err
	}

	pay.HandleChargeResult(result)

	c.payRepo.Save(context.Background(), pay)

	r := dto.NewChargeResult(0, "success", pay.Uuid().String())
	return r, nil
}
