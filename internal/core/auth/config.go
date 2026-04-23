package core_auth

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type AuthService struct {
	JWTSecret string
}

func NewAuthService() (AuthService, error) {
	var config AuthService
	if err := envconfig.Process("AUTH", &config); err != nil {
		return AuthService{}, err
	}
	return config, nil
}

func NewAuthServiceMust() AuthService {
	config, err := NewAuthService()
	if err != nil {
		panic(fmt.Errorf("couldn't get auth config: %w", err))
	}
	return config
}
