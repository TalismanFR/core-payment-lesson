package application

type TerminalRepo interface {
	GetAlias(terminalId string) (string, error)
}
