package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func HashPassword(password string) (string, error) {
	bpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bpassword), err
}

func CompareHashedPassword(password string, hPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hPassword), []byte(password))
	return err == nil
}
