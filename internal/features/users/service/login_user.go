package users_service

import core_domain "github.com/Fitray/sentinel-service/internal/core/domain"

func (s *UsersService) LoginUser(loginRequest core_domain.LoginRequest) (
	core_domain.User, error,
) {
	user, err := s.usersRepository.LoginUser(loginRequest)
	return user, err
}
