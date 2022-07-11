package dto

import (
	"diLesson/application/domain/credit_card"
	"fmt"
)

type ChargeRequest struct {
	Amount      int64
	Currency    string
	TerminalId  string
	InvoiceId   string
	Description string
	CreditCard  credit_card.CreditCard
}

func NewChargeRequest(amount int64, currency string, terminalId string, invoiceId string, description string, creditCard credit_card.CreditCard) *ChargeRequest {
	return &ChargeRequest{Amount: amount, Currency: currency, TerminalId: terminalId, InvoiceId: invoiceId, Description: description, CreditCard: creditCard}
}

func (c ChargeRequest) Valid() error {
	if err := c.CreditCard.Validate(); err != nil {
		return fmt.Errorf("credir card is invalid: %w", err)
	}
	if c.Amount < 0 {
		return fmt.Errorf("amount less than zero")
	}
	if c.Currency == "" {
		return fmt.Errorf("currency is empty")
	}
	if c.TerminalId == "" {
		return fmt.Errorf("terminalId is empty")
	}
	if c.InvoiceId == "" {
		return fmt.Errorf("invoiceId is empty")
	}
	if c.Description == "" {
		return fmt.Errorf("description is empty")
	}

	return nil
}
