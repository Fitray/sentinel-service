package users_transport_http

import (
	"net/http"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_server "github.com/Fitray/sentinel-service/internal/core/server"
)

type UsersTransport struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(userRequest core_domain.RegisterRequest) (core_domain.User, error)
	LoginUser(loginRequest core_domain.LoginRequest) (core_domain.User, error)
}

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// TODO: add auth interface here to make login possible
func NewUsersTransport(
	usersService UsersService,
) *UsersTransport {
	return &UsersTransport{
		usersService: usersService,
	}
}

func (t *UsersTransport) GetRoutes() []core_server.Route {
	return []core_server.Route{
		{
			Method:  http.MethodPost,
			Pattern: "/auth/register",
			Handler: t.CreateUser,
			Version: "v1",
		},

		{
			Method:  http.MethodGet,
			Pattern: "/auth/login",
			Handler: t.LoginUser,
			Version: "v1",
		},
	}
}
