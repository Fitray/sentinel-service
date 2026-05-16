package users_transport_http

import (
	"net/http"

	core_auth "github.com/Fitray/sentinel-service/internal/core/auth"
	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_middleware "github.com/Fitray/sentinel-service/internal/core/middleware"
	core_server "github.com/Fitray/sentinel-service/internal/core/server"
)

type UsersTransport struct {
	usersService UsersService
	authService  AuthService
}

type UsersService interface {
	CreateUser(userRequest core_domain.RegisterRequest) (core_domain.User, error)
	LoginUser(loginRequest core_domain.LoginRequest) (core_domain.User, error)
	Me(user_id string) (core_domain.User, error)
}

type AuthService interface {
	GenerateToken(user core_domain.User) (string, error)
}

// type Route struct {
// 	Method      string
// 	Pattern     string
// 	Handler     http.HandlerFunc
// }

// TODO: add auth interface here to make login possible
func NewUsersTransport(
	usersService UsersService,
	authService AuthService,
) *UsersTransport {
	return &UsersTransport{
		usersService: usersService,
		authService:  authService,
	}
}

func (t *UsersTransport) GetRoutes(
	routes []core_server.Route,
	logger core_logger.Logger,
	auth core_auth.AuthService,
) []core_server.Route {
	for _, v := range []core_server.Route{
		{
			Method:      http.MethodPost,
			Pattern:     "/auth/register",
			Handler:     t.CreateUser,
			Version:     "v1",
			Middlewares: core_middleware.NoAuthGroupMiddlewares(logger),
		},

		{
			Method:      http.MethodPost,
			Pattern:     "/auth/login",
			Handler:     t.LoginUser,
			Version:     "v1",
			Middlewares: core_middleware.NoAuthGroupMiddlewares(logger),
		},

		{
			Method:      http.MethodGet,
			Pattern:     "/auth/me",
			Handler:     t.Me,
			Version:     "v1",
			Middlewares: core_middleware.AuthGroupMiddlewares(logger, auth),
		},
	} {
		routes = append(routes, v)
	}
	return routes
}
