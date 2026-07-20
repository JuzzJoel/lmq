package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// Alphabet defines the allowed characters in the generated tokens.
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// DefaultTokenLength specifies the default length of generated tokens.
	DefaultTokenLength = 6
)

// GenerateToken creates a short, URL-friendly random string using crypto/rand.
func GenerateToken(length int) (string, error) {
	if length <= 0 {
		length = DefaultTokenLength
	}
	
	bytes := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(Alphabet)))
	
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate token: %w", err)
		}
		bytes[i] = Alphabet[num.Int64()]
	}
	
	return string(bytes), nil
}

