package core_auth

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type AuthConfig struct {
	JWTSecret string
}

func NewAuthConfig() (AuthConfig, error) {
	var config AuthConfig
	if err := envconfig.Process("AUTH", &config); err != nil {
		return AuthConfig{}, err
	}
	return config, nil
}

func NewAuthConfigMust() AuthConfig {
	config, err := NewAuthConfig()
	if err != nil {
		panic(fmt.Errorf("couldn't get auth config: %w", err))
	}
	return config
}
