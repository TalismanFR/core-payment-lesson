package application

import (
	"diLesson/application/domain/vo"
	"github.com/google/uuid"
)

type TerminalRepo interface {
	FindByUuid(terminalUuid uuid.UUID) (*vo.Terminal, error)
}
