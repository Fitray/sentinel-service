package sentinel_service_http

import "context"

type ImageryService struct {
	imageryRepositiry ImageryRepository
}

type ImageryRepository interface {
	GetImagery(ctx context.Context, city, from, to string) ([]byte, error)
}

func NewImageryService(
	imageryRepositiry ImageryRepository,
) ImageryService {
	return ImageryService{
		imageryRepositiry: imageryRepositiry,
	}
}
