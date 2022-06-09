package application

import (
	"context"
	"diLesson/application/domain"
	"github.com/google/uuid"
)

type PayRepository interface {
	Save(ctx context.Context, pay *domain.Pay) error
	Update(ctx context.Context, pay *domain.Pay) error
	FindByInvoiceID(ctx context.Context, invoiceId string) (*domain.Pay, error)
	FindByUuid(ctx context.Context, uuid uuid.UUID) (*domain.Pay, error)
}
