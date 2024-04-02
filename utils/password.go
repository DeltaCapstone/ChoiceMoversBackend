package utils

import (
	"crypto/rand"
	"encoding/base64"
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

// generateRandomCode generates a random code of the specified length
func GenerateRandomCode(length int) (string, error) {
	// Determine the number of bytes needed for the specified length
	numBytes := (length * 6) / 8 // 6 bits per base64 character

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode random bytes to base64
	randomCode := base64.URLEncoding.EncodeToString(randomBytes)

	// Truncate the code to the desired length
	if len(randomCode) > length {
		randomCode = randomCode[:length]
	}

	return randomCode, nil
}
