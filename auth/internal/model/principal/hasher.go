package principal

import "golang.org/x/crypto/bcrypt"

var global Hasher = BcryptHasher{bcrypt.DefaultCost}

func SetHasher(h Hasher) {
	global = h
}

func GetHasher() Hasher {
	return global
}

type Hasher interface {
	GenerateFromPassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hash []byte, password []byte) error
}

type BcryptHasher struct {
	cost int
}

func (b BcryptHasher) GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, b.cost)
}

func (b BcryptHasher) CompareHashAndPassword(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
