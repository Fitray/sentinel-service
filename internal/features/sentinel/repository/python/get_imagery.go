package sentinel_repository_py

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (h ImageryRepository) getCmd(ctx context.Context, city, from, to string) *exec.Cmd {
	cmd_name := fmt.Sprintf("%s/.venv/bin/python", h.Root)
	path := fmt.Sprintf("%s/internal/python/main.py", h.Root)
	return exec.CommandContext(ctx, cmd_name, path, city, from, to)
}

func (h *ImageryRepository) GetImagery(ctx context.Context, city, from, to string) ([]byte, error) {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		h.Timeout,
	)
	defer cancel()

	cmd := h.getCmd(ctxTimeout, city, from, to)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return []byte{}, fmt.Errorf("python service runtime: %w", err)
		} else {
			return []byte{}, fmt.Errorf("forbidden output from python service: %w: %w",
				core_errors.ErrBadGateway, err)
		}
	}

	return output, nil
}
