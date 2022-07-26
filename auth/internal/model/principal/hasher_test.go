package principal

import (
	"reflect"
	"testing"
)

func TestSetHasher(t *testing.T) {
	type args struct {
		h Hasher
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetHasher(tt.args.h)
		})
	}
}

func TestGetHasher(t *testing.T) {
	tests := []struct {
		name string
		want Hasher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHasher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHasher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBcryptHasher_GenerateFromPassword(t *testing.T) {
	type fields struct {
		cost int
	}
	type args struct {
		password []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BcryptHasher{
				cost: tt.fields.cost,
			}
			got, err := b.GenerateFromPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("BcryptHasher.GenerateFromPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BcryptHasher.GenerateFromPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBcryptHasher_CompareHashAndPassword(t *testing.T) {
	type fields struct {
		cost int
	}
	type args struct {
		hash     []byte
		password []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BcryptHasher{
				cost: tt.fields.cost,
			}
			if err := b.CompareHashAndPassword(tt.args.hash, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("BcryptHasher.CompareHashAndPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
