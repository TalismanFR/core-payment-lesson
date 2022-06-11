package domain

import (
	"diLesson/application/domain/vo"
	"diLesson/payment/contract/dto"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type statusCode int

const (
	StatusNew     = statusCode(0)
	StatusPending = statusCode(1)
)

type Pay struct {
	uuid          uuid.UUID
	amount        vo.Amount
	currency      vo.Currency
	invoiceId     string
	statusCode    statusCode
	status        string
	createdAt     time.Time
	transactionId string
	terminal      *vo.Terminal
}

func NewPay(uuid uuid.UUID, amount vo.Amount, currency vo.Currency, invoiceId string, statusCode statusCode, status string, createdAt time.Time, transactionId string, terminal *vo.Terminal) (*Pay, error) {
	if invoiceId == "" || status == "" || transactionId == "" {
		return nil, fmt.Errorf("invalid arguments: empty string")
	}
	return &Pay{uuid: uuid, amount: amount, currency: currency, invoiceId: invoiceId, statusCode: statusCode, status: status, createdAt: createdAt, transactionId: transactionId, terminal: terminal}, nil
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

func (p Pay) Currency() vo.Currency {
	return p.currency
}

func (p Pay) InvoiceId() string {
	return p.invoiceId
}

func (p Pay) StatusCode() int {
	return int(p.statusCode)
}

func (p Pay) Status() string {
	return p.status
}

func (p Pay) CreatedAt() time.Time {
	return p.createdAt
}

func (p Pay) TransactionId() string {
	return p.transactionId
}

func (p *Pay) IsStatusNew() bool {
	return p.statusCode == StatusNew
}

func (p *Pay) IsStatusPending() bool {
	return p.statusCode == StatusPending
}

func (p *Pay) Terminal() *vo.Terminal {
	return p.terminal
}
