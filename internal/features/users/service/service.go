package users_service

import (
	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(userRequest core_domain.RegisterRequest) (core_domain.User, error)
	LoginUser(loginRequest core_domain.LoginRequest) (core_domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
