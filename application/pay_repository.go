package application

import "diLesson/application/domain"

type PayRepository interface {
	Save(pay *domain.Pay) error
	Update(pay *domain.Pay) error
	FindByInvoiceID(invoiceId string) (*domain.Pay, error)
}
