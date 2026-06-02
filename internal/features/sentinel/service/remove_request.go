package sentinel_service

import (
	"context"
	"fmt"
	"strconv"

	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (s *ImageryService) RemoveRequestById(
	ctx context.Context, user_id string, idStr string,
) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %v: %w", err, core_errors.ErrBadRequest)
	}

	err = s.imageryRepositiry.RemoveRequestById(ctx, user_id, id)
	if err != nil {
		return err
	}

	return nil
}
