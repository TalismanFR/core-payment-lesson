package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/domain/vo"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/infrastructure/terminal"
	"diLesson/payment/bepaid"
	contract2 "diLesson/payment/contract"
	"diLesson/payment/tinkoff"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
)

func BuildDI() (err error) {
	err = container.Transient(func() (application.PayRepository, error) {
		dsn := "host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable"

		return repository.NewPayRepositoryPgsql(dsn)
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
		terminals := map[string]*vo.Terminal{
			"8242df35-e182-4448-a99d-fd6b86dd7312": vo.NewTerminal(uuid.MustParse("8242df35-e182-4448-a99d-fd6b86dd7312"), "bepaid",
				map[string]string{"login": "login", "password": "password"}),
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
