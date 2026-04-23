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

func (t *UsersTransport) LoginUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.GetLoggerFromContext(r.Context())
	rw := core_responce.NewResponceHandler(w, logger)

	var loginRequest core_domain.LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &loginRequest); err != nil {
		rw.ErrorResponce(
			fmt.Errorf("failed to decode json: %w", err),
			core_errors.GetStatusCode(err),
		)
		return
	}

	user, err := t.usersService.LoginUser(loginRequest)
	if err != nil {
		rw.ErrorResponce(
			fmt.Errorf("failed to login user: %w", err),
			core_errors.GetStatusCode(err),
		)
		return
	}

	token, err := t.authService.GenerateToken(user)
	if err != nil {
		rw.ErrorResponce(
			fmt.Errorf("failed to generate token: %w", err),
			core_errors.GetStatusCode(err),
		)
		return
	}

	responce := core_domain.LoginResponce{
		Token: token,
	}

	rw.JSONResponce(&responce, http.StatusCreated)
}
