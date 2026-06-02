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
	GetHistory(
		ctx context.Context, requestFilter core_domain.FilterRequest,
	) ([]core_domain.NewImagery, error)
	GetRequestFromID(
		ctx context.Context, user_id string, id int,
	) (core_domain.NewImagery, error)
	RemoveRequestById(
		ctx context.Context, user_id string, id int,
	) error
}

func NewImageryService(
	imageryRepositiry ImageryRepository,
) ImageryService {
	return ImageryService{
		imageryRepositiry: imageryRepositiry,
	}
}
