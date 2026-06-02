package users_service

import (
	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
)

func (s *UsersService) Me(user_id string) (
	core_domain.User, error,
) {
	user, err := s.usersRepository.GetUser("id=$1", user_id)
	if err != nil {
		return core_domain.User{}, err
	}
	return user, nil
}
