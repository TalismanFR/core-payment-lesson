package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/payment/bepaid"
	contract2 "diLesson/payment/contract"
	"github.com/golobby/container/v3"
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
