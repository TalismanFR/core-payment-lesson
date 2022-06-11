package status

import "fmt"

const (
	_                       = iota
	newCode, newStr         = iota, "NEW"
	pendingCode, pendingStr = iota, "PENDING"
)

var (
	StatusUnknown = Status{}
	StatusNew     = Status{newCode, newStr, false}
	StatusPending = Status{pendingCode, pendingStr, false}
)

type Status struct {
	code        int
	description string
	final       bool
}

func (s Status) Code() int {
	return s.code
}

func (s Status) Description() string {
	return s.description
}

func (s Status) Final() bool {
	return s.final
}

func FromString(s string) (Status, error) {

	switch s {
	case newStr:
		return StatusNew, nil
	case pendingStr:
		return StatusPending, nil
	}

	return StatusUnknown, fmt.Errorf("unknown status: %s", s)
}
