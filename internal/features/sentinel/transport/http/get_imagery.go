package sentinel_transport_http

import (
	"fmt"
	"net/http"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
)

func (h *ImageryTransport) GetImagery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.GetLoggerFromContext(ctx)
	httpHandler := core_responce.NewResponceHandler(w, logger)
	user_id := ctx.Value("user_id").(string)
	if user_id == "" {
		httpHandler.ErrorResponce(
			fmt.Errorf("failed to get user id: %w", core_errors.ErrUnauthorized),
			core_errors.GetStatusCode(core_errors.ErrUnauthorized),
		)
		return
	}

	imageryRequest := h.getDTO(r, user_id)

	output, err := h.imageryService.GetImagery(ctx, imageryRequest)
	if err != nil {
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}

	contentType := http.DetectContentType(output.Image)

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(output.Image)
}

func (h *ImageryTransport) getDTO(r *http.Request, user_id string) core_domain.ImageryRequest {
	city := r.URL.Query().Get("city")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	req := core_domain.ImageryRequest{
		City:    city,
		From:    from,
		To:      to,
		User_id: user_id,
	}
	return req
}
