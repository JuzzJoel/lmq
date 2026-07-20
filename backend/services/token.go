package services

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	// Alphabet defines the allowed characters in the generated tokens.
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// DefaultTokenLength specifies the default length of generated tokens.
	DefaultTokenLength = 6
)

// GenerateToken creates a short, URL-friendly random string using nanoid.
func GenerateToken(length int) (string, error) {
	if length <= 0 {
		length = DefaultTokenLength
	}
	token, err := gonanoid.Generate(Alphabet, length)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}
