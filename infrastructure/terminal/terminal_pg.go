package terminal

import (
	"context"
	"diLesson/application/domain/vo"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type TerminalDB struct {
	Uuid      uuid.UUID `gorm:"primary_key"`
	Alias     string
	Url       string
	CreatedAt time.Time
}

type RepoPG struct {
	db *gorm.DB
}

func NewRepoPG(dsn string) (*RepoPG, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&TerminalDB{})
	if err != nil {
		return nil, err
	}

	return &RepoPG{db}, nil
}

func (repo *RepoPG) FindByUuid(ctx context.Context, terminalUuid uuid.UUID) (*vo.Terminal, error) {

	if terminalUuid.String() == "" {
		return nil, fmt.Errorf("uuid has zero length")
	}

	tx := repo.db.WithContext(ctx)

	var t TerminalDB

	r := tx.First(&t, "uuid = ?", terminalUuid.String())

	if r.Error != nil {
		return nil, r.Error
	}

	ter := terminalFromTerminalDB(&t)

	return ter, nil
}

func terminalFromTerminalDB(t *TerminalDB) *vo.Terminal {
	return vo.NewTerminal(t.Uuid, t.Alias, map[string]string{"url": t.Url})
}
