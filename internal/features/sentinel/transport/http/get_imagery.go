package sentinel_transport_http

import (
	"context"
	"fmt"
	"net/http"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_http_request "github.com/Fitray/sentinel-service/internal/core/request"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
)

func (h *ImageryTransport) GetImageryRequest(
	r *http.Request,
) (core_domain.ImageryRequest, error) {
	var imageryRequest core_domain.ImageryRequest
	err := core_http_request.DecodeAndValidateRequest(r, &imageryRequest)
	if err != nil {
		return core_domain.ImageryRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return imageryRequest, nil
}

func (h *ImageryTransport) GetContent(
	ctx context.Context, imageryRequest core_domain.ImageryRequest,
) (core_domain.ImageryResponce, error) {
	output, err := h.imageryService.GetImagery(ctx, imageryRequest)
	if err != nil {
		return core_domain.ImageryResponce{}, err
	}
	return output, nil
}

func (h *ImageryTransport) GetPreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.GetLoggerFromContext(ctx)
	httpHandler := core_responce.NewResponceHandler(w, logger)

	user_id, ok := ctx.Value("user_id").(string)
	if !ok || user_id == "" {
		err := fmt.Errorf("failed to get user id: %w", core_errors.ErrUnauthorized)
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}

	imageryRequest, err := h.GetImageryRequest(r)
	if err != nil {
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}
	imageryRequest.User_id = user_id
	imageryRequest.OutputFormat = "png"

	output, err := h.GetContent(ctx, imageryRequest)
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

func (h *ImageryTransport) DownloadImagery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.GetLoggerFromContext(ctx)
	httpHandler := core_responce.NewResponceHandler(w, logger)

	user_id, ok := ctx.Value("user_id").(string)
	if !ok || user_id == "" {
		err := fmt.Errorf("failed to get user id: %w", core_errors.ErrUnauthorized)
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}

	imageryRequest, err := h.GetImageryRequest(r)
	if err != nil {
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}
	imageryRequest.User_id = user_id

	output, err := h.GetContent(ctx, imageryRequest)
	if err != nil {
		httpHandler.ErrorResponce(
			err,
			core_errors.GetStatusCode(err),
		)
		return
	}

	contentType := http.DetectContentType(output.Image)
	fileName := fmt.Sprintf("%s %s-%s.tif", imageryRequest.City, imageryRequest.From, imageryRequest.To)
	header := fmt.Sprintf("attachment; filename=\"%s\"", fileName)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", header)
	w.WriteHeader(http.StatusOK)
	w.Write(output.Image)
}
