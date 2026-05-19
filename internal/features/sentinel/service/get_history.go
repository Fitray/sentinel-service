package sentinel_service

import (
	"context"
	"fmt"
	"strconv"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (s *ImageryService) GetHistory(
	ctx context.Context, requestFilter core_domain.FilterRequest,
) ([]core_domain.NewImagery, error) {
	if requestFilter.From.After(requestFilter.To) {
		return []core_domain.NewImagery{},
			fmt.Errorf("date from can't be larger than date to: %w",
				core_errors.ErrBadRequest,
			)
	}
	reqHistory, err := s.imageryRepositiry.GetHistory(ctx, requestFilter)
	if err != nil {
		return []core_domain.NewImagery{}, err
	}
	return reqHistory, nil
}

func (s *ImageryService) GetRequestFromID(
	ctx context.Context, user_id string, idStr string,
) (core_domain.NewImagery, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return core_domain.NewImagery{},
			fmt.Errorf("invalid id: %v: %w",
				core_errors.ErrBadRequest, err)
	}

	resp, err := s.imageryRepositiry.GetRequestFromID(ctx, user_id, id)
	if err != nil {
		return core_domain.NewImagery{}, err
	}

	return resp, nil
}
