package principal

import (
	"testing"
)

func testPrincipal() *Principal {
	return &Principal{
		Email:    "example@ya.ru",
		Password: "password",
		Role:     "user",
	}
}

func TestPrincipal_Validate(t *testing.T) {
	tests := []struct {
		name  string
		u     func() *Principal
		valid bool
	}{
		{
			name: "valid_1",
			u: func() *Principal {
				return testPrincipal()
			},
			valid: true,
		},
		{
			name: "invalid_1",
			u: func() *Principal {
				u := testPrincipal()
				u.Email = ""
				u.Password = ""
				u.Role = ""

				return u
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.u().Validate()
			if tc.valid && err != nil {
				t.Fatalf("error isn't nil:\nER: nil\nAR: %s", err)
			}

			if !tc.valid && err == nil {
				t.Fatalf("error is nil:\nER: not nil error\nAR: nil")
			}
		})
	}
}

func TestPrincipal_Sanitize(t *testing.T) {
}

func TestPrincipal_HashPassword(t *testing.T) {

}

func TestPrincipal_PasswordMatches(t *testing.T) {

}
