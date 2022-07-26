package principal

import (
	"encoding"
	"encoding/json"
	validate "github.com/go-playground/validator/v10"
)

var validator = validate.New()

var _ encoding.BinaryMarshaler = (*Principal)(nil)
var _ encoding.BinaryUnmarshaler = (*Principal)(nil)

type Principal struct {
	Email          string `validate:"required,email"`
	Password       string `validate:"min=8,max=16"`
	Role           string
	HashedPassword []byte
}

func (p *Principal) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)

}

func (p *Principal) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

//TODO: add role validation
func (p *Principal) Validate() error {
	return validator.Struct(p)
}

func (p *Principal) Sanitize() {
	p.Password = ""
}

func (p *Principal) HashPassword() error {
	b, err := global.GenerateFromPassword([]byte(p.Password))
	if err != nil {
		return err
	}

	p.HashedPassword = b
	return nil
}

func (p *Principal) PasswordMatches(password string) bool {
	return global.CompareHashAndPassword(p.HashedPassword, []byte(password)) == nil
}
