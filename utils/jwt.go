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
func GenerateToken(adminID uint, version int) (string, time.Time, error) {
	// Calculate expiration time
	expiryMinutes := config.Config.JWT.ExpiryMinutes
	expiryTime := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	// Create claims
	claims := jwt.MapClaims{
		"admin_id":      adminID,
		"token_version": version,
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
func ParseToken(tokenStr string) (uint, int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, 0, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return 0, 0, errors.New("missing expiration time")
		}

		if time.Now().Unix() > int64(exp) {
			return 0, 0, errors.New("token expired")
		}

		// Extract admin ID
		id, ok1 := claims["admin_id"].(float64)
		if !ok1 {
			return 0, 0, errors.New("invalid admin ID")
		}

		// Extract token version
		ver, ok2 := claims["token_version"].(float64)
		if !ok2 {
			return 0, 0, errors.New("invalid token version")
		}

		return uint(id), int(ver), nil
	}

	return 0, 0, errors.New("invalid token")
}
