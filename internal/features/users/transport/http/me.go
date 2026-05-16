package users_transport_http

import (
	"fmt"
	"net/http"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
)

func (t *UsersTransport) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.GetLoggerFromContext(ctx)
	httpHandler := core_responce.NewResponceHandler(w, logger)

	user_id, ok := ctx.Value("user_id").(string)
	if !ok || user_id == "" {
		httpHandler.ErrorResponce(
			fmt.Errorf("failed to get user id: %w", core_errors.ErrUnauthorized),
			core_errors.GetStatusCode(core_errors.ErrUnauthorized),
		)
		return
	}

	user, err := t.usersService.Me(user_id)
	if err != nil {
		httpHandler.ErrorResponce(
			fmt.Errorf("failed to get user: %w"),
			core_errors.GetStatusCode(err),
		)
		return
	}

	userResponce := core_domain.NewRegisterResonce(user)
	httpHandler.JSONResponce(userResponce, http.StatusOK)
}
