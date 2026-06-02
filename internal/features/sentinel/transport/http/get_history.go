package sentinel_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_responce "github.com/Fitray/sentinel-service/internal/core/responce"
	"github.com/go-chi/chi/v5"
)

func (t *ImageryTransport) GetHistory(
	w http.ResponseWriter, r *http.Request,
) {
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

	requestFilter, err := t.getFilterDTO(r, user_id)
	if err != nil {
		httpHandler.ErrorResponce(err, http.StatusBadRequest)
		return
	}

	reqHistory, err := t.imageryService.GetHistory(ctx, requestFilter)
	if err != nil {
		httpHandler.ErrorResponce(err, core_errors.GetStatusCode(err))
		return
	}

	httpHandler.JSONResponce(reqHistory, http.StatusOK)
}

func (t *ImageryTransport) GetRequestFromID(w http.ResponseWriter, r *http.Request) {
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
	resp, err := t.imageryService.GetRequestFromID(ctx, user_id, idStr)
	if err != nil {
		httpHandler.ErrorResponce(
			err, core_errors.GetStatusCode(err),
		)
		return
	}

	httpHandler.JSONResponce(resp, http.StatusOK)
}

func (h *ImageryTransport) getFilterDTO(
	r *http.Request, user_id string,
) (core_domain.FilterRequest, error) {
	city := r.URL.Query().Get("city")
	from_str := r.URL.Query().Get("from")
	to_str := r.URL.Query().Get("to")
	var (
		from time.Time
		to   time.Time
		err  error
	)

	if from_str != "" {
		from, err = time.Parse("2006-01-02", from_str)
		if err != nil {
			return core_domain.FilterRequest{}, fmt.Errorf(
				"failed to parse 'from': %w: %v",
				core_errors.ErrBadRequest,
				err,
			)
		}
	}

	if to_str != "" {
		to, err = time.Parse("2006-01-02", to_str)
		if err != nil {
			return core_domain.FilterRequest{}, fmt.Errorf(
				"failed to parse 'to': %w: %v",
				core_errors.ErrBadRequest,
				err,
			)
		}
	}

	return core_domain.FilterRequest{
		City:    city,
		From:    from,
		To:      to,
		User_id: user_id,
	}, nil
}
