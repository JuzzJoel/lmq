package services

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	// Alphabet defines the allowed characters in the generated tokens.
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// TokenLength specifies the default length of generated tokens.
	TokenLength = 6
)

// GenerateToken creates a short, URL-friendly random string using nanoid.
func GenerateToken() (string, error) {
	token, err := gonanoid.Generate(Alphabet, TokenLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}
