package core_middleware

import (
	"net/http"

	core_auth "github.com/Fitray/sentinel-service/internal/core/auth"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	chi "github.com/go-chi/chi/v5"
)

type Middleware func(next http.Handler) http.Handler

func AuthGroupMiddlewares(
	logger core_logger.Logger,
	auth core_auth.AuthService,
) chi.Middlewares {
	return chi.Middlewares{
		RequestID(),
		Logger(logger),
		Panic(),
		NewRequest(),
		AuthMiddleware(auth),
	}
}

func NoAuthGroupMiddlewares(
	logger core_logger.Logger,
) chi.Middlewares {
	return chi.Middlewares{
		RequestID(),
		Logger(logger),
		Panic(),
		NewRequest(),
	}
}
