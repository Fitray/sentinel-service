package users_repository_postgres

import (
	core_postgres_pool "github.com/Fitray/sentinel-service/internal/core/postgres/pool"
)

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	connPool core_postgres_pool.ConnectionPool,
) *UsersRepository {
	return &UsersRepository{
		pool: connPool,
	}
}
