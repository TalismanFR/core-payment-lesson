package principal

import (
	validate "github.com/go-playground/validator/v10"
)

var validator = validate.New()

type Principal struct {
	Email          string `validate:"required,email"`
	Password       string `validate:"min=8,max=16"`
	Role           string
	HashedPassword []byte
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
