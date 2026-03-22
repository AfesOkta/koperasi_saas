package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
)

// GenerateToken creates a new JWT access token.
func GenerateToken(userID, orgID, roleID uint, roleVersion int, email, secret string, expirationHours int) (string, error) {
	claims := middleware.Claims{
		UserID:         userID,
		OrganizationID: orgID,
		Email:          email,
		RoleID:         roleID,
		RoleVersion:    roleVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken creates a new JWT refresh token.
func GenerateRefreshToken(userID uint, secret string, refreshHours int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshHours) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyRefreshToken validates the refresh token and returns the user ID.
func VerifyRefreshToken(tokenString, secret string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		var id uint
		fmt.Sscanf(claims.Subject, "%d", &id)
		return id, nil
	}
	return 0, fmt.Errorf("invalid refresh token")
}
