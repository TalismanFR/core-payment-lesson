package terminal

import (
	"context"
	"diLesson/application/domain/terminal"
	"fmt"
	"github.com/google/uuid"
)

type TerminalRepoInMemory struct {
	terminals map[string]*terminal.Terminal
}

func NewTerminalRepoInMemory(terminals map[string]*terminal.Terminal) *TerminalRepoInMemory {
	return &TerminalRepoInMemory{terminals: terminals}
}

func (t *TerminalRepoInMemory) FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*terminal.Terminal, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		a, ok := t.terminals[terminalUuid.String()]
		if !ok {
			return nil, fmt.Errorf("no such terminalUuid")
		}

		return a, nil
	}
}
