package users_service

import (
	"fmt"
	"net/mail"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) CreateUser(
	userRequest core_domain.RegisterRequest,
) (core_domain.User, error) {
	if _, err := mail.ParseAddress(userRequest.Email); err != nil {
		return core_domain.User{}, fmt.Errorf("email check: %v: %w", err, core_errors.ErrBadRequest)
	}

	if len([]rune(userRequest.Password)) < 6 {
		return core_domain.User{}, fmt.Errorf("short password: %w", core_errors.ErrBadRequest)
	}

	if len([]rune(userRequest.Name)) < 3 {
		return core_domain.User{}, fmt.Errorf("short name: %w", core_errors.ErrBadRequest)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("creating user: %w", err)
	}
	userRequest.Password = string(hash)

	user, err := s.usersRepository.CreateUser(userRequest)
	if err != nil {
		return user, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}
