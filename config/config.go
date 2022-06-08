package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/infrastructure/terminal"
	"diLesson/payment/bepaid"
	contract2 "diLesson/payment/contract"
	"diLesson/payment/tinkoff"
	"github.com/golobby/container/v3"
)

func BuildDI() (err error) {
	err = container.Transient(func() application.PayRepository {
		return &repository.PayRepositoryPgsql{}
	})

	err = container.NamedTransient("bepaid", func() (contract2.VendorCharge, error) {
		s := bepaid.Charge{}

		return &s, nil
	})

	err = container.NamedTransient("tinkoff", func() (contract2.VendorCharge, error) {
		s := tinkoff.Charge{}

		return &s, nil
	})

	err = container.Transient(func() application.TerminalRepo {
		terminals := map[string]string{
			"terminalId1": "bepaid",
			"terminalId2": "tinkoff",
			"terminalId3": "bepaid",
		}

		return terminal.NewTerminalRepoInMemory(terminals)
	})

	err = container.Transient(func() (contract.Charge, error) {
		var ch service.Charge

		err := container.Fill(&ch)

		return &ch, err
	})

	return err
}
