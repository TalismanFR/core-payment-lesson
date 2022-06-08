package terminal

import "fmt"

type TerminalRepoInMemory struct {
	terminals map[string]string
}

func NewTerminalRepoInMemory(terminals map[string]string) *TerminalRepoInMemory {
	return &TerminalRepoInMemory{terminals: terminals}
}

func (t *TerminalRepoInMemory) GetAlias(terminalId string) (string, error) {
	a, ok := t.terminals[terminalId]
	if !ok {
		return "", fmt.Errorf("no such terminalId")
	}

	return a, nil
}
