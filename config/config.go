package config

import (
	"github.com/golobby/container/v3"
	"payservice-core/application"
	"payservice-core/application/contract"
	"payservice-core/application/service"
	"payservice-core/infrastructure/repository"
	"payservice-core/payment/bepaid"
	contract2 "payservice-core/payment/contract"
)

func BuildDI() (err error) {
    err = container.Transient(func() application.PayRepository {
        return &repository.PayRepositoryPgsql{}
    })

    err = container.NamedTransient("bepaid_charge", func() (contract2.VendorCharge, error) {
        s := bepaid.Charge{}

        return &s, nil
    })

    err = container.Transient(func() (contract.Charge, error) {
        var ch service.Charge

        err := container.Fill(&ch)

        return &ch, err
    })

    return err
}
