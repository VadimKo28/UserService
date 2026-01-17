package hash

import (
	"crypto/sha1"
	"fmt"
)

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
	  return "", err
	}

	hashedPassword := hash.Sum([]byte(h.salt))

	return fmt.Sprintf("%x", hashedPassword), nil
}