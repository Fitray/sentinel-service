package sentinel_service

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

var (
	correctBands = []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B10", "B11", "B12"}
)

func (h *ImageryService) GetImagery(
	ctx context.Context, imageryRequest core_domain.ImageryRequest,
) (core_domain.ImageryResponce, error) {
	if imageryRequest.City == "" {
		return core_domain.ImageryResponce{}, fmt.Errorf("forbidden city %w", core_errors.ErrInvalidArg)
	}

	if imageryRequest.Id != "" {
		if _, err := strconv.Atoi(imageryRequest.Id); err != nil {
			return core_domain.ImageryResponce{},
				fmt.Errorf("invalid id type: %v: %w", err, core_errors.ErrBadRequest)
		}
	}

	bands := strings.Split(imageryRequest.Bands, ",")
	for _, band := range bands {
		if !slices.Contains(correctBands, band) {
			return core_domain.ImageryResponce{}, fmt.Errorf("invalid layer: %s: %w", band, core_errors.ErrBadRequest)
		}
	}

	if imageryRequest.Dimensions <= 0 || imageryRequest.Dimensions > 2048 {
		return core_domain.ImageryResponce{}, fmt.Errorf("invalid dimensions: %d: %w", imageryRequest.Dimensions, core_errors.ErrBadRequest)
	}

	if imageryRequest.Scale <= 0 || imageryRequest.Scale >= 100 {
		return core_domain.ImageryResponce{}, fmt.Errorf("invalid scale: %d: %w", imageryRequest.Scale, core_errors.ErrBadRequest)
	}

	if imageryRequest.OutputFormat != "png" && imageryRequest.OutputFormat != "geotiff" {
		return core_domain.ImageryResponce{}, fmt.Errorf("invalid output format: %s: %w", imageryRequest.OutputFormat, core_errors.ErrBadRequest)
	}

	output, err := h.imageryRepositiry.GetImagery(ctx, imageryRequest)
	return output, err
}
