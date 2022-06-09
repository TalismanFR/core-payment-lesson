package repository

import (
	"context"
	"diLesson/application/domain"
	"diLesson/application/domain/vo"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Pay struct {
	Uuid          uuid.UUID `gorm:"primaryKey"`
	Amount        int
	Currency      string
	InvoiceId     string
	StatusCode    int
	Status        string
	CreatedAt     time.Time
	TransactionId string
}

type PayRepositoryPgsql struct {
	db *gorm.DB
}

func NewPayRepositoryPgsql(dsn string) (*PayRepositoryPgsql, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Pay{})
	if err != nil {
		return nil, err
	}

	return &PayRepositoryPgsql{db}, nil
}

func (repo *PayRepositoryPgsql) Save(ctx context.Context, pay *domain.Pay) error {
	tx := repo.db.WithContext(ctx)
	r := tx.Create(payFromDomainPay(pay))

	return r.Error
}

func (repo *PayRepositoryPgsql) Update(ctx context.Context, pay *domain.Pay) error {
	tx := repo.db.WithContext(ctx)

	p := payFromDomainPay(pay)

	r := tx.First(&Pay{Uuid: pay.Uuid()}).Updates(p)

	if r.RowsAffected == 0 {
		return fmt.Errorf("no record with UUID: %v", p.Uuid)
	}

	return r.Error
}

func (repo *PayRepositoryPgsql) FindByInvoiceID(ctx context.Context, invoiceId string) (*domain.Pay, error) {

	if invoiceId == "" {
		return nil, fmt.Errorf("empty invoiceId")
	}

	tx := repo.db.WithContext(ctx)

	var pay *Pay

	r := tx.First(pay, "invoiceId = ?", invoiceId)

	if r.Error != nil {
		return nil, r.Error
	}

	return domainPayFromPay(pay), nil
}

func (repo *PayRepositoryPgsql) FindByUuid(ctx context.Context, uuid uuid.UUID) (*domain.Pay, error) {

	if uuid.String() == "" {
		return nil, fmt.Errorf("uuid with zero length")
	}

	tx := repo.db.WithContext(ctx)

	var pay *Pay

	r := tx.First(pay, "id = ?", uuid.String())

	if r.Error != nil {
		return nil, r.Error
	}

	return domainPayFromPay(pay), nil
}

func payFromDomainPay(pay *domain.Pay) *Pay {
	return &Pay{
		pay.Uuid(),
		int(pay.Amount()),
		string(pay.Currency()),
		pay.InvoiceId(),
		pay.StatusCode(),
		pay.Status(),
		pay.CreatedAt(),
		pay.TransactionId(),
	}
}

func domainPayFromPay(pay *Pay) *domain.Pay {

	code := domain.StatusCodeOK

	switch pay.StatusCode {
	case int(domain.StatusCodeOK):
		code = domain.StatusCodeOK
	}

	p, _ := domain.NewPay(
		pay.Uuid,
		vo.Amount(pay.Amount),
		vo.Currency(pay.Currency),
		pay.InvoiceId,
		code,
		pay.Status,
		pay.CreatedAt,
		pay.TransactionId,
	)

	return p
}
