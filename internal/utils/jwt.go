package utils

import (
	"errors"
	"fmt"
	"time"

	"TinderTrip-Backend/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Provider  string `json:"provider"`
	ExpiresAt int64  `json:"expires_at"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID, email, provider string) (string, error) {
	cfg := config.AppConfig.JWT

	// Create claims
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		Provider:  provider,
		ExpiresAt: time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "TinderTrip-Backend",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.AppConfig.JWT

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	// Check if token is expired
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiration
func RefreshToken(tokenString string) (string, error) {
	// Validate current token
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Generate new token
	return GenerateToken(claims.UserID, claims.Email, claims.Provider)
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[7:], nil
}

// GetTokenExpiration returns the token expiration time
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(claims.ExpiresAt, 0), nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(tokenString string) bool {
	expiration, err := GetTokenExpiration(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(expiration)
}
