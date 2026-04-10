package core_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr    string        `envconfig:"ADDR" required="true"`
	Timeout time.Duration `envconfig:"TIMEOUT" required="true"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func MustConfig() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get HTTP server config: %w", err))
	}
	return config
}
