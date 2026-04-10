package core_middleware

import (
	"net/http"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	"github.com/go-chi/chi"
)

type Middleware func(next http.Handler) http.Handler

func GetMiddlewares(logger core_logger.Logger) chi.Middlewares {
	return chi.Middlewares{
		RequestID(),
		Logger(logger),
		Panic(),
		NewRequest(),
	}
}
