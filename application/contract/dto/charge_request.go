package dto

import "fmt"

type ChargeRequest struct {
	Amount      int
	Currency    string
	TerminalId  string
	InvoiceId   string
	Description string
}

func (c ChargeRequest) Valid() error {
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
