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
	StatusCodeOK = statusCode(0)
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
}

func NewPay(uuid uuid.UUID, amount vo.Amount, currency vo.Currency, invoiceId string, statusCode statusCode, status string, createdAt time.Time, transactionId string) (*Pay, error) {
	if invoiceId == "" || status == "" || transactionId == "" {
		return nil, fmt.Errorf("invalid arguments: empty string")
	}
	return &Pay{uuid: uuid, amount: amount, currency: currency, invoiceId: invoiceId, statusCode: statusCode, status: status, createdAt: createdAt, transactionId: transactionId}, nil
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
