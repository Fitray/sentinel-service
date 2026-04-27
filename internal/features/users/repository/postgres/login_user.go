package users_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (r *UsersRepository) LoginUser(loginRequest core_domain.LoginRequest) (
	core_domain.User, error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		r.pool.GetTimeout(),
	)
	defer cancel()

	query := `
	SELECT id, full_name, email, password_hash, created_at, updated_at FROM app.users
	WHERE email=$1
	`

	row := r.pool.QueryRow(ctx, query, loginRequest.Email)

	var user core_domain.User

	err := row.Scan(
		&user.ID, &user.Name, &user.Email,
		&user.Password, &user.Created_at, &user.Updated_at,
	)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("invalid credentials: %v: %w",
			err, core_errors.ErrUnauthorized,
		)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(loginRequest.Password),
	); err != nil {
		return core_domain.User{}, fmt.Errorf("invalid credentials: %v: %w",
			err, core_errors.ErrUnauthorized)
	}

	return user, nil
}
