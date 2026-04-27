package sentinel_service

import (
	"context"
	"fmt"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (h *ImageryService) GetImagery(
	ctx context.Context, imageryRequest core_domain.ImageryRequest,
) (core_domain.ImageryResponce, error) {
	if imageryRequest.City == "" {
		return core_domain.ImageryResponce{}, fmt.Errorf("forbidden city %w", core_errors.ErrInvalidArg)
	}

	output, err := h.imageryRepositiry.GetImagery(ctx, imageryRequest)
	return output, err
}
