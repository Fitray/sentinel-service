package sentinel_service

import (
	"context"
	"fmt"

	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (h *ImageryService) GetImagery(ctx context.Context, city, from, to string) ([]byte, error) {
	if city == "" {
		return []byte{}, fmt.Errorf("forbidden city %w", core_errors.ErrInvalidArg)
	}

	output, err := h.imageryRepositiry.GetImagery(ctx, city, from, to)
	return output, err
}
