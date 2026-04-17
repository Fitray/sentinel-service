package sentinel_repository_http

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ImageryRepository struct {
	Root    string        `envconfig:"PROJECT_ROOT" required="true"`
	Timeout time.Duration `envconfig:"PYTHON_TIMEOUT" required="true"`
}

func NewImageryRepository() (ImageryRepository, error) {
	var conf ImageryRepository
	if err := envconfig.Process("", &conf); err != nil {
		return ImageryRepository{}, err
	}
	return conf, nil
}

func NewImageryRepositoryMust() ImageryRepository {
	conf, err := NewImageryRepository()
	if err != nil {
		panic(fmt.Errorf("failed to get sentinel repository config: %w", err))
	}
	return conf
}
