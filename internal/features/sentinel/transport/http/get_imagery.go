package sentinel_transport_http

import (
	"net/http"

	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
)

func (h *ImageryTransport) GetImagery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.GetLoggerFromContext(ctx)
	httpHandler := core_responce.NewResponceHandler(w, logger)
	city := r.URL.Query().Get("city")

	output, err := h.imageryService.GetImagery(ctx, city)
	if err != nil {
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
	}

	contentType := http.DetectContentType(output)

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
