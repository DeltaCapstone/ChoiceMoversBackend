package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(bytes), nil
}
