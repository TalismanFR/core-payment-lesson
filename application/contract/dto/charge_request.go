package dto

import (
	"diLesson/application/domain/vo"
	"fmt"
)

type ChargeRequest struct {
	Amount      int
	Currency    string
	TerminalId  string
	InvoiceId   string
	Description string
	CreditCard  *vo.CreditCard
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
