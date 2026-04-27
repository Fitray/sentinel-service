package sentinel_service

import (
	"context"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
)

type ImageryService struct {
	imageryRepositiry ImageryRepository
}

type ImageryRepository interface {
	GetImagery(
		ctx context.Context, imageryRequest core_domain.ImageryRequest,
	) (core_domain.ImageryResponce, error)
}

func NewImageryService(
	imageryRepositiry ImageryRepository,
) ImageryService {
	return ImageryService{
		imageryRepositiry: imageryRepositiry,
	}
}
