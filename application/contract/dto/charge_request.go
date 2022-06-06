package dto

type ChargeRequest struct {
	Amount      int
	TerminalId  string
	InvoiceId   string
	Description string
}
