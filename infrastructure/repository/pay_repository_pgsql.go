package repository

import (
    "fmt"
    "payservice-core/application/domain"
)

type PayRepositoryPgsql struct {
}

func (p PayRepositoryPgsql) Save(pay *domain.Pay) error {
    fmt.Println("save pay ", pay.Uuid())
    return nil
}

func (p PayRepositoryPgsql) Update(pay *domain.Pay) error {
    //TODO implement me
    panic("implement me")
}

func (p PayRepositoryPgsql) FindByInvoiceID(invoiceId string) (*domain.Pay, error) {
    //TODO implement me
    panic("implement me")
}
