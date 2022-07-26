package helpers

import (
	"fmt"
	"testing"
)

func ExampleParseField() {
	input := `
keyspace_misses:123
keyspace_hits:456`
	field := "keyspace_misses:"
	chars := "\n\r\n\r"

	hits, _ := ParseField([]byte(input), []byte(field), chars)

	fmt.Printf("%[1]T, %[1]v", hits)

	// Output:
	// int64, 123
}

func TestParseField(t *testing.T) {
	type args struct {
		input []byte
		field []byte
		chars string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseField(tt.args.input, tt.args.field, tt.args.chars)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseField() = %v, want %v", got, tt.want)
			}
		})
	}
}
