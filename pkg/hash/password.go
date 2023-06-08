package hash

import (
	"crypto/sha256"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type sha256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *sha256Hasher {
	return &sha256Hasher{salt: salt}
}

func (h *sha256Hasher) Hash(password string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
