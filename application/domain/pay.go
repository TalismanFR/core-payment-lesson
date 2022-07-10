package domain

import (
	"diLesson/payment/contract/dto"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Pay struct {
	uuid        uuid.UUID
	amount      int
	currency    string //todo to struc
	invoiceId   string
	statusCode  int //todo to struct
	status      string
	createdAt   time.Time
	transaction string
	terminal    *Terminal
}

func (p *Pay) HandleChargeResult(result *dto.VendorChargeResult) {
	result.IsFailed()
	//p.status=400
}

func (p Pay) Uuid() uuid.UUID {
	return p.uuid
}

func (p Pay) Amount() int {
	return p.amount
}

func (p Pay) Currency() string {
	return p.currency
}

func (p Pay) InvoiceId() string {
	return p.invoiceId
}

func (p Pay) StatusCode() int {
	return p.statusCode
}

func (p Pay) Status() string {
	return p.status
}

func (p Pay) CreatedAt() time.Time {
	return p.createdAt
}

func (p Pay) Transaction() string {
	return p.transaction
}

func (p *Pay) Terminal() *Terminal {
	return p.terminal
}

func NewPay(uuid uuid.UUID, amount int, currency string, invoiceId string, statusCode int, status string, createdAt time.Time, transaction string) *Pay {
	return &Pay{uuid: uuid, amount: amount, currency: currency, invoiceId: invoiceId, statusCode: statusCode, status: status, createdAt: createdAt, transaction: transaction}
}
