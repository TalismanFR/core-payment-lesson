package terminal

import (
	"diLesson/application/domain/vo"
	"fmt"
	"github.com/google/uuid"
)

type TerminalRepoInMemory struct {
	terminals map[string]*vo.Terminal
}

func NewTerminalRepoInMemory(terminals map[string]*vo.Terminal) *TerminalRepoInMemory {
	return &TerminalRepoInMemory{terminals: terminals}
}

func (t *TerminalRepoInMemory) FindByUuid(terminalUuid uuid.UUID) (*vo.Terminal, error) {
	a, ok := t.terminals[terminalUuid.String()]
	if !ok {
		return nil, fmt.Errorf("no such terminalUuid")
	}

	return a, nil
}
