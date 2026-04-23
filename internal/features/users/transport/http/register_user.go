package users_transport_http

import (
	"fmt"
	"net/http"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_http_request "github.com/Fitray/sentinel-service/internal/core/request"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
)

func (t *UsersTransport) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.GetLoggerFromContext(r.Context())
	responceHandler := core_responce.NewResponceHandler(w, logger)

	var userRequest core_domain.RegisterRequest
	err := core_http_request.DecodeAndValidateRequest(r, &userRequest)
	if err != nil {
		responceHandler.ErrorResponce(
			fmt.Errorf("failed to decode request: %w", err),
			http.StatusBadRequest,
		)
		return
	}

	user, err := t.usersService.CreateUser(userRequest)
	if err != nil {
		responceHandler.ErrorResponce(
			fmt.Errorf("failed to create user: %w"),
			core_errors.GetStatusCode(err),
		)
		return
	}

	userResponce := core_domain.NewRegisterResonce(user)
	responceHandler.JSONResponce(userResponce, http.StatusCreated)
}
