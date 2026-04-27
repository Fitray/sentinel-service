package sentinel_repository_py

import (
	"fmt"
	"time"

	core_postgres_pool "github.com/Fitray/sentinel-service/internal/core/postgres/pool"
	"github.com/kelseyhightower/envconfig"
)

type ImageryRepository struct {
	Root    string        `envconfig:"PROJECT_ROOT" required="true"`
	Timeout time.Duration `envconfig:"PYTHON_TIMEOUT" required="true"`
	Pool    core_postgres_pool.Pool
}

func NewImageryRepository(
	pool core_postgres_pool.ConnectionPool,
) (ImageryRepository, error) {
	var conf ImageryRepository
	if err := envconfig.Process("", &conf); err != nil {
		return ImageryRepository{}, err
	}
	conf.Pool = pool
	return conf, nil
}

func NewImageryRepositoryMust(
	pool core_postgres_pool.ConnectionPool,
) ImageryRepository {
	conf, err := NewImageryRepository(pool)
	if err != nil {
		panic(fmt.Errorf("failed to get sentinel repository config: %w", err))
	}
	return conf
}
