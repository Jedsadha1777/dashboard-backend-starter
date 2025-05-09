package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// GenerateRandomPassword generates a random password with the specified length
// ensuring it contains at least one character from each character class
func GenerateRandomPassword(length int) string {
	if length < 8 {
		length = 8 // Enforce minimum length
	}

	// Ensure we have at least one character from each class
	password := make([]byte, length)

	// Start with one character from each required class
	addRandomChar(password, 0, lowerChars)
	addRandomChar(password, 1, upperChars)
	addRandomChar(password, 2, digitChars)
	addRandomChar(password, 3, specialChars)

	// Fill the rest with random characters from all classes
	allChars := lowerChars + upperChars + digitChars + specialChars
	for i := 4; i < length; i++ {
		addRandomChar(password, i, allChars)
	}

	// Shuffle the password to randomize positions
	shuffleBytes(password)

	return string(password)
}

// addRandomChar adds a random character from the given character set to the password at the specified index
func addRandomChar(password []byte, index int, charSet string) {
	max := big.NewInt(int64(len(charSet)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to a safe character if there's an error
		password[index] = 'A'
		return
	}
	password[index] = charSet[n.Int64()]
}

// shuffleBytes randomly shuffles a byte slice
func shuffleBytes(bytes []byte) {
	for i := len(bytes) - 1; i > 0; i-- {
		// Get random index
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue // Skip shuffle step on error
		}

		// Swap elements
		bytes[i], bytes[j.Int64()] = bytes[j.Int64()], bytes[i]
	}
}

// IsStrongPassword checks if a password meets strength requirements
func IsStrongPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	checks := []struct {
		pattern string
		message string
	}{
		{lowerChars, "Password must contain at least one lowercase letter"},
		{upperChars, "Password must contain at least one uppercase letter"},
		{digitChars, "Password must contain at least one digit"},
		{specialChars, "Password must contain at least one special character"},
	}

	for _, check := range checks {
		if !containsAny(password, check.pattern) {
			return false, check.message
		}
	}

	return true, ""
}

// containsAny checks if the string contains any character from the given character set
func containsAny(s string, chars string) bool {
	return strings.IndexAny(s, chars) >= 0
}
