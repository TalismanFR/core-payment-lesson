package repository

import (
	"context"
	"diLesson/application/domain"
	"diLesson/application/domain/currency"
	"diLesson/application/domain/status"
	"diLesson/application/domain/vo"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Pay struct {
	Uuid          uuid.UUID `gorm:"primaryKey"`
	Amount        uint64
	Currency      string
	InvoiceId     string
	Status        string
	CreatedAt     time.Time
	TransactionId string
	TerminalId    string
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

	p, err := domainPayFromPay(pay)
	if err != nil {
		return nil, fmt.Errorf("couldn't create domain.Pay from Pay: %w", err)
	}

	return p, nil
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

	p, err := domainPayFromPay(pay)
	if err != nil {
		return nil, fmt.Errorf("couldn't create domain.Pay from Pay: %w", err)
	}

	return p, nil

}

func payFromDomainPay(pay *domain.Pay) *Pay {
	return &Pay{
		pay.Uuid(),
		uint64(pay.Amount()),
		pay.Currency().String(),
		pay.InvoiceId(),
		pay.Status().Description(),
		pay.CreatedAt(),
		pay.TransactionId(),
		pay.Terminal().Uuid().String(),
	}
}

func domainPayFromPay(pay *Pay) (*domain.Pay, error) {

	stat, err := status.FromString(pay.Status)
	if err != nil {
		return nil, err
	}

	cur, err := currency.FromString(pay.Currency)
	if err != nil {
		return nil, err
	}

	terminalUuid, err := uuid.Parse(pay.TerminalId)
	if err != nil {
		return nil, err
	}

	p, _ := domain.PayFull(
		pay.Uuid,
		vo.Amount(pay.Amount),
		cur,
		"empty description",
		pay.InvoiceId,
		stat,
		pay.CreatedAt,
		pay.TransactionId,
		vo.NewTerminal(terminalUuid, "", nil),
		nil,
	)

	return p, nil
}
