package core_middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	core_auth "github.com/Fitray/sentinel-service/internal/core/auth"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
	"github.com/google/uuid"
)

const (
	requestID = "X-Request-ID"
	loggerKey = "logger"
)

func AuthMiddleware(authConfig core_auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			userId, err := authConfig.GetUserIDfromHeader(authHeader)

			if err != nil {
				logger := core_logger.GetLoggerFromContext(r.Context())
				responceHandler := core_responce.NewResponceHandler(w, logger)
				responceHandler.ErrorResponce(
					fmt.Errorf("failed to authorize user: %w", err),
					core_errors.GetStatusCode(err),
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				"user_id", userId,
			)

			rw := r.WithContext(ctx)
			next.ServeHTTP(w, rw)
		})
	}
}

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
			log := core_logger.GetLoggerFromContext(r.Context())
			defer func() {
				if p := recover(); p != nil {
					responseHandler := core_responce.NewResponceHandler(w, log)
					responseHandler.PanicResponce("a panic occurred during server runtime", p)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func NewRequest() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.GetLoggerFromContext(r.Context())
			reqId := r.Header.Get(requestID)
			responseHandler := core_responce.NewResponceHandler(w, log)

			next.ServeHTTP(responseHandler, r)

			log.Info("new request",
				slog.String(requestID, reqId),
				slog.String("Method", r.Method),
				slog.Int("statusCode", responseHandler.StatusCode),
				slog.String("Pattern", r.Pattern),
			)
		})
	}
}
