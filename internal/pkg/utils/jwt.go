package utils

import (
	"errors"
	"netfilessys/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateResetToken generates a password reset token (valid for 1 hour)
func GenerateResetToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret + "_reset"))
}

// ValidateResetToken validates a reset token and returns user ID
func ValidateResetToken(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret + "_reset"), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired token")
	}

	return claims.UserID, nil
}

// ParseTokenAllowExpired parses a token even if it's expired (useful for refresh)
func ParseTokenAllowExpired(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		// Check if error is due to expiration
		if ve, ok := err.(*jwt.ValidationError); ok {
			// If the error includes expiration, allow it for refresh
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired but otherwise valid if we have claims with UserID
				if claims.UserID > 0 {
					return claims, nil
				}
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateToken validates a token and returns user ID
func ValidateToken(tokenString string) (uint, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ValidateTokenAllowExpired validates a token even if expired (useful for refresh)
func ValidateTokenAllowExpired(tokenString string) (uint, error) {
	claims, err := ParseTokenAllowExpired(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
