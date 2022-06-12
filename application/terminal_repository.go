package application

import (
	"context"
	"diLesson/application/domain/vo"
	"github.com/google/uuid"
)

type TerminalRepo interface {
	FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*vo.Terminal, error)
}
