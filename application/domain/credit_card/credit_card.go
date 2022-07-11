package credit_card

import (
	"fmt"
)

type CreditCard struct {

	// номер карты, длина - от 12 до 19 цифр
	number string `json:"number"`

	//3-х или 4-х цифровой код безопасности (CVC2, CVV2 или CID, в зависимости от бренда карты).
	//Может быть отправлен вместе с параметром token и bePaid доставит банку-эквайеру данные карты с CVC2/CVV2/CID
	verificationValue string `json:"verification_value"`

	//имя владельца карты. Максимальная длина: 32 символа
	holder string `json:"holder"`

	//месяц окончания срока действия карты, представленный двумя цифрами (например, 01)
	expMonth ExpirationMonth `json:"exp_month"`

	//год срока окончания действия карты, представленный четырьмя цифрами (например, 2007)
	expYear string `json:"exp_year"`

	//опционально
	//вместо 5 параметров выше можно отправить токен карты, который был получен в ответе первой оплаты.
	//Если используется токен карты, то необходимо обязательно указывать параметр additional_data.contract
	token string `json:"token,omitempty"`

	//опционально
	//если значение параметра true, bePaid не выполняет 3-D Secure проверку.
	//Это полезно если вы, например, не хотите чтобы клиент проходил 3-D Secure проверку снова.
	//Уточните у службы поддержки, можете ли вы использовать этот параметр.
	skipThreeDSecureVerification bool `json:"skip_three_d_secure_verification"`
}

func (cc *CreditCard) Number() string {
	return cc.number
}

func (cc *CreditCard) VerificationValue() string {
	return cc.verificationValue
}

func (cc *CreditCard) Holder() string {
	return cc.holder
}

func (cc *CreditCard) ExpMonth() ExpirationMonth {
	return cc.expMonth
}

func (cc *CreditCard) ExpYear() string {
	return cc.expYear
}

func (cc *CreditCard) Token() string {
	return cc.token
}

func (cc *CreditCard) SkipThreeDSecureVerification() bool {
	return cc.skipThreeDSecureVerification
}

func NewCreditCard(number string, verificationValue string, holder string, expMonth ExpirationMonth, expYear string) *CreditCard {
	return &CreditCard{
		number:            number,
		verificationValue: verificationValue,
		holder:            holder,
		expMonth:          expMonth,
		expYear:           expYear,
	}
}

func (cc *CreditCard) Validate() error {
	if len(cc.number) < 12 || len(cc.number) > 19 {
		return fmt.Errorf("number length should be between 12 and 19 (both including)")
	}
	if len(cc.verificationValue) != 3 && len(cc.verificationValue) != 4 {
		return fmt.Errorf("verificationCode length should equal to 3 or 4")
	}
	if len(cc.holder) > 32 {
		return fmt.Errorf("holder length greater than 32")
	}
	if cc.expMonth == UnknownExpMonth {
		return fmt.Errorf("invalid expiration month")
	}
	if len(cc.expYear) != 4 {
		return fmt.Errorf("expYear length should equal 4")
	}

	return nil
}

func NewCreditCardWithToken(token string) *CreditCard {
	return &CreditCard{token: token}
}

func (cc *CreditCard) WithSkip3DSVerification(skipThreeDSecureVerification bool) *CreditCard {
	cc.skipThreeDSecureVerification = skipThreeDSecureVerification
	return cc
}

type ExpirationMonth struct {
	i int
	s string
}

func (e ExpirationMonth) Int() int {
	return e.i
}

func (e ExpirationMonth) String() string {
	return e.s
}

func FromString(s string) (ExpirationMonth, error) {
	switch s {
	case janStr:
		return January, nil
	case febStr:
		return February, nil
	case marStr:
		return March, nil
	case aprStr:
		return April, nil
	case mayStr:
		return May, nil
	case junStr:
		return June, nil
	case julStr:
		return July, nil
	case augStr:
		return August, nil
	case sepStr:
		return September, nil
	case octStr:
		return October, nil
	case novStr:
		return November, nil
	case decStr:
		return December, nil
	}

	return UnknownExpMonth, fmt.Errorf("unknown expiration month: %q", s)
}

func FromInt(i int) (ExpirationMonth, error) {
	switch i {
	case janInt:
		return January, nil
	case febInt:
		return February, nil
	case marInt:
		return March, nil
	case aprInt:
		return April, nil
	case mayInt:
		return May, nil
	case junInt:
		return June, nil
	case julInt:
		return July, nil
	case augInt:
		return August, nil
	case sepInt:
		return September, nil
	case octInt:
		return October, nil
	case novInt:
		return November, nil
	case decInt:
		return December, nil
	}

	return UnknownExpMonth, fmt.Errorf("unknown expiration month: %q", i)
}

const (
	_              = iota
	janInt, janStr = iota, "01"
	febInt, febStr = iota, "02"
	marInt, marStr = iota, "03"
	aprInt, aprStr = iota, "04"
	mayInt, mayStr = iota, "05"
	junInt, junStr = iota, "06"
	julInt, julStr = iota, "07"
	augInt, augStr = iota, "08"
	sepInt, sepStr = iota, "09"
	octInt, octStr = iota, "10"
	novInt, novStr = iota, "11"
	decInt, decStr = iota, "12"
)

var (
	UnknownExpMonth = ExpirationMonth{0, ""}
	January         = ExpirationMonth{janInt, janStr}
	February        = ExpirationMonth{febInt, febStr}
	March           = ExpirationMonth{marInt, marStr}
	April           = ExpirationMonth{aprInt, aprStr}
	May             = ExpirationMonth{mayInt, mayStr}
	June            = ExpirationMonth{junInt, junStr}
	July            = ExpirationMonth{julInt, julStr}
	August          = ExpirationMonth{augInt, augStr}
	September       = ExpirationMonth{sepInt, sepStr}
	October         = ExpirationMonth{octInt, octStr}
	November        = ExpirationMonth{novInt, novStr}
	December        = ExpirationMonth{decInt, decStr}
)
