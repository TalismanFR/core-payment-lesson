package service

import (
	"context"
	"diLesson/application"
	"diLesson/application/contract/dto"
	"diLesson/application/domain"
	"diLesson/application/domain/currency"
	"diLesson/application/domain/vo"
	"diLesson/payment/contract"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
)

const (
	chargeSuffix = "_charge"
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

//func (c Charge) Update(id uuid.UUID, request dto.ChargeRequest) error {
//	p, err := domain.NewPay(id, vo.Amount(request.Amount), vo.RUB, request.InvoiceId, domain.StatusCodeOK, vo.StatusNew, time.Now(), "transactionId")
//	if err != nil {
//		return err
//	}
//	return c.payRepo.Update(context.Background(), p)
//}

func (c Charge) Charge(request dto.ChargeRequest) (*dto.ChargeResult, error) {

	if err := request.Valid(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	cur, err := currency.FromString(request.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid requst: %w", err)
	}

	// TODO: when to generate and when to acquire transactionId
	uuidTerminal, err := uuid.Parse(request.TerminalId)
	if err != nil {
		return nil, err
	}

	terminal, err := c.terminalRepo.FindByUuid(uuidTerminal)

	pay, err := domain.NewPay(uuid.New(), vo.Amount(request.Amount), cur, request.InvoiceId, terminal)
	if err != nil {
		return nil, err
	}

	var vendor contract.VendorCharge
	err = container.NamedResolve(&vendor, terminal.Alias()+chargeSuffix)

	result, err := vendor.Charge(pay)
	if err != nil {
		return nil, err
	}

	pay.HandleChargeResult(result)

	err = c.payRepo.Save(context.Background(), pay)
	if err != nil {
		return nil, err
	}

	r := dto.NewChargeResult(0, "success", pay.Uuid().String())
	return r, nil
}
