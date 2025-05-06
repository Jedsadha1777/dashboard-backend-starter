package utils

import (
	"dashboard-starter/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

// InitJWT initializes the JWT signing key
func InitJWT() error {
	// Get JWT secret from config
	secret := config.Config.JWT.Secret
	if secret == "" {
		return errors.New("JWT_SECRET not set in environment")
	}

	jwtKey = []byte(secret)
	return nil
}

// GenerateToken creates a new JWT token for the specified admin
func GenerateToken(userID uint, userType string, version int) (string, time.Time, error) {
	// Calculate expiration time
	expiryMinutes := config.Config.JWT.ExpiryMinutes
	expiryTime := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	// Create claims
	claims := jwt.MapClaims{
		"user_id":       userID,
		"user_type":     userType, // Add user type
		"token_version": version,
		"token_type":    "access", // Specify it's an access token
		"exp":           expiryTime.Unix(),
		"iat":           time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

// ParseToken validates a JWT token and returns the admin ID and token version
func ParseToken(tokenStr string) (uint, string, int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, "", 0, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return 0, "", 0, errors.New("missing expiration time")
		}

		if time.Now().Unix() > int64(exp) {
			return 0, "", 0, errors.New("token expired")
		}

		// Extract user ID
		id, ok1 := claims["user_id"].(float64)
		if !ok1 {
			return 0, "", 0, errors.New("invalid user ID")
		}

		// Extract user type
		userType, ok2 := claims["user_type"].(string)
		if !ok2 {
			return 0, "", 0, errors.New("invalid user type")
		}

		// Extract token version
		ver, ok3 := claims["token_version"].(float64)
		if !ok3 {
			return 0, "", 0, errors.New("invalid token version")
		}

		return uint(id), userType, int(ver), nil
	}

	return 0, "", 0, errors.New("invalid token")
}

// GenerateRefreshToken creates a refresh token for the specified user
func GenerateRefreshToken(userID uint, userType string) (string, time.Time, error) {
	// Calculate expiration time (1 year)
	expiryDays := 365 // 1 year
	expiryTime := time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)

	// Create claims
	claims := jwt.MapClaims{
		"user_id":    userID,
		"user_type":  userType,  // Add user type
		"token_type": "refresh", // Specify it's a refresh token
		"exp":        expiryTime.Unix(),
		"iat":        time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

// ParseRefreshToken validates a refresh token and returns the user information
func ParseRefreshToken(tokenStr string) (uint, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return 0, "", errors.New("missing expiration time")
		}

		if time.Now().Unix() > int64(exp) {
			return 0, "", errors.New("token expired")
		}

		// Extract user ID
		id, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", errors.New("invalid user ID")
		}

		// Extract user type
		userType, ok := claims["user_type"].(string)
		if !ok {
			return 0, "", errors.New("invalid user type")
		}

		// Verify token type
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != "refresh" {
			return 0, "", errors.New("invalid token type")
		}

		return uint(id), userType, nil
	}

	return 0, "", errors.New("invalid token")
}
