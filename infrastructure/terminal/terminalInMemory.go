package terminal

import (
	"context"
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

func (t *TerminalRepoInMemory) FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*vo.Terminal, error) {
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
