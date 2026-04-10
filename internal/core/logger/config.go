package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Path string `envconfig:"PATH" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func MustConfig() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Sprintf("an error occured during getting logger config: %s", err.Error()))
	}
	return config
}
