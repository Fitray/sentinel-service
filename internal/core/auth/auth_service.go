package core_auth

import (
	"fmt"
	"time"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (c *AuthConfig) GenerateToken(user core_domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(c.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}
