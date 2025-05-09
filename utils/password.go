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
	minLength    = 12 // เพิ่มความยาวขั้นต่ำของรหัสผ่านเป็น 12 ตัวอักษร
)

// GenerateRandomPassword generates a random password with the specified length
// ensuring it contains at least one character from each character class
func GenerateRandomPassword(length int) string {
	if length < minLength {
		length = minLength // Enforce minimum length of 12
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
// ใช้ crypto/rand เพื่อความปลอดภัย
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
	if len(password) < minLength {
		return false, "Password must be at least 12 characters long"
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

	// คำตรวจสอบรหัสผ่านที่ใช้บ่อยหรือง่ายเกินไป
	commonPasswords := []string{
		"password", "123456", "qwerty", "admin", "welcome",
		"password123", "admin123", "welcome123",
	}

	lowerPassword := strings.ToLower(password)
	for _, commonPwd := range commonPasswords {
		if strings.Contains(lowerPassword, commonPwd) {
			return false, "Password contains common or easily guessable phrases"
		}
	}

	// ตรวจสอบลำดับตัวอักษรหรือตัวเลขที่ต่อเนื่องกัน
	for i := 0; i < len(password)-2; i++ {
		if isConsecutive(password[i : i+3]) {
			return false, "Password contains sequential characters or numbers"
		}
	}

	return true, ""
}

// isConsecutive ตรวจสอบว่าตัวอักษรหรือตัวเลข 3 ตัวเรียงกันหรือไม่ เช่น "123", "abc"
func isConsecutive(s string) bool {
	// ตรวจสอบตัวเลขเรียงกัน
	if isDigit(s[0]) && isDigit(s[1]) && isDigit(s[2]) {
		if (int(s[1]) == int(s[0])+1) && (int(s[2]) == int(s[1])+1) {
			return true
		}
		if (int(s[1]) == int(s[0])-1) && (int(s[2]) == int(s[1])-1) {
			return true
		}
	}

	// ตรวจสอบตัวอักษรเรียงกัน
	if isLetter(s[0]) && isLetter(s[1]) && isLetter(s[2]) {
		lower := strings.ToLower(s)
		if (lower[1] == lower[0]+1) && (lower[2] == lower[1]+1) {
			return true
		}
		if (lower[1] == lower[0]-1) && (lower[2] == lower[1]-1) {
			return true
		}
	}

	return false
}

// isDigit ตรวจสอบว่าเป็นตัวเลขหรือไม่
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLetter ตรวจสอบว่าเป็นตัวอักษรหรือไม่
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// containsAny checks if the string contains any character from the given character set
func containsAny(s string, chars string) bool {
	return strings.IndexAny(s, chars) >= 0
}
