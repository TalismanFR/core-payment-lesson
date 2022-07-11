package application

import (
	"context"
	"diLesson/application/domain/terminal"
	"github.com/google/uuid"
)

type TerminalRepo interface {
	FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*terminal.Terminal, error)
}
