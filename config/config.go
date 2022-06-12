package config

import (
	"diLesson/application"
	"diLesson/application/contract"
	"diLesson/application/service"
	"diLesson/infrastructure/repository"
	"diLesson/infrastructure/secrets"
	"diLesson/infrastructure/terminal"
	"diLesson/payment/bepaid"
	contractPayment "diLesson/payment/contract"
	"diLesson/payment/tinkoff"
	"github.com/golobby/container/v3"
)

func BuildDI() (err error) {

	err = container.Transient(func() (application.SecretsService, error) {
		//TODO: move host and mountPath to config file
		v, e := secrets.NewVaultService("http://127.0.0.1:8200", "terminals")
		if e != nil {
			return nil, e
		}

		if v.Validate() != nil {
			return nil, e
		}

		return v, nil
	})

	err = container.Transient(func() (application.PayRepository, error) {
		//TODO: move dsn to config file
		dsn := "host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable"

		return repository.NewPayRepositoryPgsql(dsn)
	})

	err = container.NamedTransient("bepaid_charge", func() (contractPayment.VendorCharge, error) {
		s := bepaid.NewCharge("bepaid")

		return s, nil
	})

	err = container.NamedTransient("tinkoff_charge", func() (contractPayment.VendorCharge, error) {
		s := tinkoff.Charge{}

		return &s, nil
	})

	err = container.Transient(func() (application.TerminalRepo, error) {
		//TODO: move dsn to config file
		dsn := "host=localhost user=payservice password=payservice dbname=payservice-db port=5432 sslmode=disable"

		return terminal.NewRepoPG(dsn)
		//terminals := map[string]*vo.Terminal{
		//	"8242df35-e182-4448-a99d-fd6b86dd7312": vo.NewTerminal(uuid.MustParse("8242df35-e182-4448-a99d-fd6b86dd7312"), "bepaid",
		//		map[string]string{"login": "login", "password": "password"}),
		//}
		//
		//return terminal.NewTerminalRepoInMemory(terminals), nil
	})

	err = container.Transient(func() (contract.Charge, error) {
		var ch service.Charge

		err := container.Fill(&ch)

		return &ch, err
	})

	return err
}
