package sentinel_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
	"github.com/go-chi/chi/v5"
)

func (t *ImageryTransport) RemoveRequestById(w http.ResponseWriter, r *http.Request) {
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

	idStr := chi.URLParam(r, "id")
	err := t.imageryService.RemoveRequestById(ctx, user_id, idStr)
	if err != nil {
		httpHandler.ErrorResponce(
			err, core_errors.GetStatusCode(err),
		)
		return
	}

	httpHandler.JSONResponce(
		map[string]interface{}{"message": "Request removed successfully"},
		http.StatusOK,
	)
}
