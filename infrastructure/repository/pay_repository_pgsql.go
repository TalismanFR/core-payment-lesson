package repository

import (
	"context"
	"diLesson/application/domain"
	"fmt"
)

type PayRepositoryPgsql struct {
}

func (p PayRepositoryPgsql) Save(ctx context.Context, pay *domain.Pay) error {
	fmt.Println("saved to db: " + pay.Uuid().String())
	return nil
}

func (p PayRepositoryPgsql) Update(ctx context.Context, pay *domain.Pay) error {
	//TODO implement me
	panic("implement me")
}

func (p PayRepositoryPgsql) FindByInvoiceID(ctx context.Context, invoiceId string) (*domain.Pay, error) {
	//TODO implement me
	panic("implement me")
}
