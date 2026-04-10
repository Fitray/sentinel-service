package core_middleware

import (
	"context"
	"log/slog"
	"net/http"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
	"github.com/google/uuid"
)

const (
	requestID = "X-Request-ID"
	loggerKey = "logger"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := r.Header.Get(requestID)
			if reqId == "" {
				r.Header.Set(requestID, uuid.NewString())
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(logger core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				loggerKey, logger,
			)
			req := r.WithContext(ctx)

			next.ServeHTTP(w, req)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.GetLoggerFromContext(r.Context().Value(loggerKey))
			defer func() {
				if p := recover(); p != nil {
					responseHandler := core_responce.NewResponceHandler(w)
					err := responseHandler.PanicResponce("a panic occurred during server runtime", p)
					log.Error(err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func NewRequest() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.GetLoggerFromContext(r.Context().Value(loggerKey))
			reqId := r.Header.Get(requestID)
			responseHandler := core_responce.NewResponceHandler(w)

			next.ServeHTTP(responseHandler, r)

			log.Info("new request",
				slog.String(requestID, reqId),
				slog.String("Method", r.Method),
				slog.Int("statusCode", responseHandler.StatusCode),
			)
		})
	}
}
