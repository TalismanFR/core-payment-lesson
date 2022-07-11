package currency

import "fmt"

var (
	Unknown = Currency{""}
	USD     = Currency{"USD"}
	RUB     = Currency{"RUB"}
	BYN     = Currency{"BYN"}
	UAH     = Currency{"UAH"}
)

type Currency struct {
	c string
}

func (cur Currency) String() string {
	return cur.c
}

func FromString(s string) (Currency, error) {
	switch s {
	case USD.String():
		return USD, nil
	case RUB.String():
		return RUB, nil
	case BYN.String():
		return BYN, nil
	case UAH.String():
		return UAH, nil
	}

	return Unknown, fmt.Errorf("unknown currency: %s", s)
}
