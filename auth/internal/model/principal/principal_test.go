package principal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrincipal_Validate(t *testing.T) {
	type fields struct {
		Email          string
		Password       string
		Role           string
		HashedPassword []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Principal{
				Email:          tt.fields.Email,
				Password:       tt.fields.Password,
				Role:           tt.fields.Role,
				HashedPassword: tt.fields.HashedPassword,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Principal.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrincipal_Sanitize(t *testing.T) {
	type fields struct {
		Email          string
		Password       string
		Role           string
		HashedPassword []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Principal{
				Email:          tt.fields.Email,
				Password:       tt.fields.Password,
				Role:           tt.fields.Role,
				HashedPassword: tt.fields.HashedPassword,
			}
			p.Sanitize()
		})
	}
}

func TestPrincipal_HashPassword(t *testing.T) {
	type fields struct {
		Email          string
		Password       string
		Role           string
		HashedPassword []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Principal{
				Email:          tt.fields.Email,
				Password:       tt.fields.Password,
				Role:           tt.fields.Role,
				HashedPassword: tt.fields.HashedPassword,
			}
			if err := p.HashPassword(); (err != nil) != tt.wantErr {
				t.Errorf("Principal.HashPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrincipal_PasswordMatches(t *testing.T) {
	type fields struct {
		Email          string
		Password       string
		Role           string
		HashedPassword []byte
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Principal{
				Email:          tt.fields.Email,
				Password:       tt.fields.Password,
				Role:           tt.fields.Role,
				HashedPassword: tt.fields.HashedPassword,
			}
			if got := p.PasswordMatches(tt.args.password); got != tt.want {
				t.Errorf("Principal.PasswordMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	p := Principal{
		Email:          "timez",
		Password:       "",
		Role:           "user",
		HashedPassword: []byte("password123"),
	}
	d, err := p.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	p1 := Principal{}
	err = p1.UnmarshalBinary(d)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, p, p1)
}
