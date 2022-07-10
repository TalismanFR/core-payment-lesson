package bepaid

import (
	"diLesson/application/domain"
	"github.com/dantedenis/bepaid/service/vo"
	"testing"
)

var card = &vo.CreditCard{
	Number:            "123123123",
	VerificationValue: "",
	Holder:            "",
	ExpMonth:          "",
	ExpYear:           "",
	Token:             "",
}

var pay = &domain.Pay{}

// какие конкретно зависимости мокать?

func TestNCharge(t *testing.T) {

}
