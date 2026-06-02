package users_repository_postgres

import (
	"fmt"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (r *UsersRepository) LoginUser(loginRequest core_domain.LoginRequest) (
	core_domain.User, error,
) {
	user, err := r.GetUser("email=$1", loginRequest.Email)
	if err != nil {
		return core_domain.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(loginRequest.Password),
	); err != nil {
		return core_domain.User{}, fmt.Errorf("invalid credentials: %v: %w",
			err, core_errors.ErrUnauthorized)
	}

	return user, nil
}
