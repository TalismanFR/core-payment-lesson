package service

import (
	"diLesson/application"
	"diLesson/application/contract/dto"
	"diLesson/application/domain"
	contract2 "diLesson/payment/contract"
	"github.com/golobby/container/v3"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Charge struct {
	repository application.PayRepository `container:"type"`
}

func (c Charge) Charge(request dto.ChargeRequest) (*dto.ChargeResult, error) {

	pay := domain.NewPay(uuid.NewV4(), request.Amount, "RUB", request.InvoiceId, 0, "new", time.Now(), "")
	//todo получение информации о вендоре

	var service contract2.VendorCharge
	container.NamedResolve(&service, "bepaid_charge")

	result, err := service.Charge(pay)
	if err != nil {
		return nil, err
	}

	pay.HandleChargeResult(result)

	c.repository.Save(pay)

	r := dto.NewChargeResult(0, "success", pay.Uuid().String())
	return r, nil
}
