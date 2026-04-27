package users_repository_postgres

import (
	"context"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	userRequest core_domain.RegisterRequest,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		r.pool.GetTimeout(),
	)
	defer cancel()

	query := `
	INSERT INTO app.users (full_name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, full_name, email, password_hash, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		userRequest.Name,
		userRequest.Email,
		userRequest.Password,
	)

	var user core_domain.User

	if err := row.Scan(
		&user.ID, &user.Name, &user.Email,
		&user.Password, &user.Created_at, &user.Updated_at,
	); err != nil {
		return core_domain.User{}, err
	}

	return user, nil
}
