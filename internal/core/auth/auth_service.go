package core_auth

import (
	"fmt"
	"strings"
	"time"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

func (c *AuthService) GenerateToken(user core_domain.User) (string, error) {
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

func (c *AuthService) GetUserIDfromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" || tokenString == authHeader {
		return "", fmt.Errorf("invalid authorization token format")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("forbidden signing method: %v", t.Header["alg"])
		}
		return []byte(c.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired token: %v: %w", err, core_errors.ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims: %v: %w", err, core_errors.ErrUnauthorized)
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims: %v: %w", err, core_errors.ErrUnauthorized)
	}

	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return "", fmt.Errorf("expired token: %v: %w", err, core_errors.ErrUnauthorized)
		}
	}

	return userId, nil
}
