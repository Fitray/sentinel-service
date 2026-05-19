package sentinel_repository

import (
	"context"
	"fmt"
)

func (h *ImageryRepository) RemoveRequestById(
	ctx context.Context, user_id string, id int,
) error {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		h.Timeout,
	)
	defer cancel()

	query := `
	DELETE FROM app.requests WHERE id=$1 AND user_id=$2
	`

	_, err := h.Pool.Exec(ctxTimeout, query, id, user_id)
	if err != nil {
		return fmt.Errorf(
			"failed to delete request: %w", err,
		)
	}

	return nil
}
