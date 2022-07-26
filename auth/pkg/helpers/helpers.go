package helpers

import (
	"bytes"
	"errors"
	"strconv"
)

// parseField cuts integer between field name and one of the characters from chars
func ParseField(input []byte, field []byte, chars string) (int64, error) {

	i := bytes.Index(input, field)
	if i == -1 {
		return -1, errors.New("")
	}
	j := bytes.IndexAny(input[i:], chars)
	if j == -1 {
		return -1, errors.New("")
	}
	v, err := strconv.Atoi(string(input[i+len(field) : i+j]))
	if err != nil {
		return -1, errors.New("")
	}

	return int64(v), nil
}
