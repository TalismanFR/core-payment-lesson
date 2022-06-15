package service

import (
    "github.com/golobby/container/v3"
    uuid "github.com/satori/go.uuid"
    "payservice-core/application"
    "payservice-core/application/contract/dto"
    "payservice-core/application/domain"
    contract2 "payservice-core/payment/contract"
    "time"
)

type Charge struct {
    repository application.PayRepository `container:"type"`
}

func (c Charge) Charge(request dto.ChargeRequest) (*dto.ChargeResult, error) {

    pay := domain.NewPay(uuid.NewV4(), request.Amount, "RUB", request.InvoiceId, 0, "new", time.Now(), "")
    //todo получение информации о вендоре

    var service contract2.VendorCharge
    if err := container.NamedResolve(&service, "bepaid_charge"); err != nil {
        return nil, err
    }

    result, err := service.Charge(pay)
    if err != nil {
        return nil, err
    }

    pay.HandleChargeResult(result)

    if err := c.repository.Save(pay); err != nil {
        return nil, err
    }

    r := dto.NewChargeResult(0, "success", pay.Uuid().String())
    return r, nil
}
