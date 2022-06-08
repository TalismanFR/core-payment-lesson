package application

import (
	"context"
	"diLesson/application/domain"
)

type PayRepository interface {
	Save(ctx context.Context, pay *domain.Pay) error
	Update(ctx context.Context, pay *domain.Pay) error
	FindByInvoiceID(ctx context.Context, invoiceId string) (*domain.Pay, error)
}
