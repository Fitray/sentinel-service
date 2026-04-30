package sentinel_repository

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (h ImageryRepository) getCmd(
	ctx context.Context,
	imageryRequest core_domain.ImageryRequest,
) *exec.Cmd {
	cmd_name := fmt.Sprintf("%s/.venv/bin/python", h.Root)
	path := fmt.Sprintf("%s/internal/python/main.py", h.Root)
	return exec.CommandContext(
		ctx, cmd_name, path,
		imageryRequest.City, imageryRequest.From, imageryRequest.To,
	)
}

func (h ImageryRepository) AddNewUserRequest(
	imageryRequest core_domain.ImageryRequest,
) (
	core_domain.NewImagery, error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		h.Pool.GetTimeout(),
	)
	defer cancel()

	query := `
	INSERT INTO app.requests (user_id, city, date_from, date_to)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, city, date_from, date_to, created_at, updated_at
	`

	row := h.Pool.QueryRow(
		ctx,
		query,
		imageryRequest.User_id,
		imageryRequest.City,
		imageryRequest.From,
		imageryRequest.To,
	)

	var newImagery core_domain.NewImagery

	if err := row.Scan(
		&newImagery.Id, &newImagery.User_id, &newImagery.City,
		&newImagery.From, &newImagery.To, &newImagery.Created_at,
		&newImagery.Updated_at,
	); err != nil {
		return core_domain.NewImagery{}, fmt.Errorf("failed to add request: %w", err)
	}

	return newImagery, nil
}

func (h *ImageryRepository) GetImagery(
	ctx context.Context, imageryRequest core_domain.ImageryRequest,
) (core_domain.ImageryResponce, error) {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		h.Timeout,
	)
	defer cancel()

	cmd := h.getCmd(ctxTimeout, imageryRequest)

	output, err := cmd.CombinedOutput()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return core_domain.ImageryResponce{}, fmt.Errorf("python service runtime: %w", err)
		} else {
			return core_domain.ImageryResponce{}, fmt.Errorf("forbidden output from python service: %w: %w",
				core_errors.ErrBadGateway, err)
		}
	}

	newUser, err := h.AddNewUserRequest(imageryRequest)
	if err != nil {
		return core_domain.ImageryResponce{}, fmt.Errorf("postgres err: %w", err)
	}

	responce := core_domain.ImageryResponce{
		Image:      output,
		NewImagery: newUser,
	}

	return responce, nil
}
