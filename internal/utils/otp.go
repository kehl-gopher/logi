package utils

import (
	"crypto/rand"
	"encoding/base64"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func Tokens() (string, error) {
	var n int = 6
	return generateSecureToken(n)
}

func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b)[:n], nil
}

func NanoId() (string, error) {
	id, err := gonanoid.New(4)
	return id, err
}
