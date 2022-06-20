package vo

import "github.com/google/uuid"

type Terminal struct {
	uuid             uuid.UUID
	alias            string
	additionalParams map[string]string
}

func (t Terminal) Uuid() uuid.UUID {
	return t.uuid
}

func (t Terminal) Alias() string {
	return t.alias
}

func (t Terminal) AdditionalParams() map[string]string {
	return t.additionalParams
}

func NewTerminal(uuid uuid.UUID, alias string, additionalParams map[string]string) *Terminal {
	return &Terminal{uuid: uuid, alias: alias, additionalParams: additionalParams}
}