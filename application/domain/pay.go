package domain

import (
	"diLesson/application/domain/currency"
	"diLesson/application/domain/status"
	"diLesson/application/domain/vo"
	"diLesson/payment/contract/dto"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Pay struct {
	uuid          uuid.UUID
	amount        vo.Amount
	currency      currency.Currency
	description   string
	invoiceId     string
	status        status.Status
	createdAt     time.Time
	transactionId string
	terminal      *vo.Terminal
	creditCard    *vo.CreditCard
}

func NewPay(uuid uuid.UUID, amount vo.Amount, currency currency.Currency, description string, invoiceId string, terminal *vo.Terminal, cc *vo.CreditCard) (*Pay, error) {
	if currency.String() == "" || invoiceId == "" {
		return nil, fmt.Errorf("invalid arguments: empty string")
	}
	return &Pay{
		uuid:          uuid,
		amount:        amount,
		currency:      currency,
		description:   description,
		invoiceId:     invoiceId,
		status:        status.StatusNew,
		createdAt:     time.Now(),
		transactionId: "",
		terminal:      terminal,
		creditCard:    cc,
	}, nil
}

func PayFull(uuid uuid.UUID, amount vo.Amount, currency currency.Currency, description string, invoiceId string, status status.Status, createdAt time.Time, transactionId string, terminal *vo.Terminal, cc *vo.CreditCard) (*Pay, error) {

	p, err := NewPay(uuid, amount, currency, description, invoiceId, terminal, cc)
	if err != nil {
		return nil, err
	}

	p.status = status
	p.createdAt = createdAt
	p.transactionId = transactionId

	return p, nil
}

func (p *Pay) HandleChargeResult(result *dto.VendorChargeResult) {
	result.IsFailed()
	//p.status=400
}

func (p Pay) Uuid() uuid.UUID {
	return p.uuid
}

func (p Pay) Amount() vo.Amount {
	return p.amount
}

func (p Pay) Currency() currency.Currency {
	return p.currency
}

func (p Pay) Description() string {
	return p.description
}

func (p Pay) InvoiceId() string {
	return p.invoiceId
}

func (p Pay) Status() status.Status {
	return p.status
}

func (p Pay) CreatedAt() time.Time {
	return p.createdAt
}

func (p Pay) TransactionId() string {
	return p.transactionId
}

func (p *Pay) IsStatusNew() bool {
	return p.status == status.StatusNew
}

func (p *Pay) IsStatusPending() bool {
	return p.status == status.StatusPending
}

func (p *Pay) Terminal() *vo.Terminal {
	return p.terminal
}

func (p *Pay) CreditCard() *vo.CreditCard {
	return p.creditCard
}
