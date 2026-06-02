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
	cmdName := fmt.Sprintf("%s/.venv/bin/python", h.Root)
	path := fmt.Sprintf("%s/internal/python/main.py", h.Root)

	return exec.CommandContext(
		ctx,
		cmdName,
		path,
		fmt.Sprintf("city=%s", imageryRequest.City),
		fmt.Sprintf("date_from=%s", imageryRequest.From),
		fmt.Sprintf("date_to=%s", imageryRequest.To),
		fmt.Sprintf("bands=%s", imageryRequest.Bands),
		fmt.Sprintf("dimensions=%d", imageryRequest.Dimensions),
		fmt.Sprintf("scale=%d", imageryRequest.Scale),
		fmt.Sprintf("output_format=%s", imageryRequest.OutputFormat),
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
	INSERT INTO app.requests (
		user_id,
		city,
		date_from,
		date_to,
		bands,
		scale,
		dimensions
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)

	ON CONFLICT (
		user_id,
		city,
		date_from,
		date_to,
		bands,
		scale,
		dimensions
	)

	DO UPDATE SET
		updated_at = NOW()

	RETURNING
		id,
		user_id,
		city,
		date_from,
		date_to,
		bands,
		scale,
		dimensions,
		created_at,
		updated_at
	`

	row := h.Pool.QueryRow(
		ctx,
		query,
		imageryRequest.User_id,
		imageryRequest.City,
		imageryRequest.From,
		imageryRequest.To,
		imageryRequest.Bands,
		imageryRequest.Scale,
		imageryRequest.Dimensions,
	)

	var newImagery core_domain.NewImagery

	if err := row.Scan(
		&newImagery.Id, &newImagery.User_id, &newImagery.City,
		&newImagery.From, &newImagery.To, &newImagery.Bands, &newImagery.Scale, &newImagery.Dimensions,
		&newImagery.Created_at, &newImagery.Updated_at,
	); err != nil {
		return core_domain.NewImagery{}, fmt.Errorf("failed to add request: %w", err)
	}

	return newImagery, nil
}

func (h ImageryRepository) GetUserRequest(
	ctx context.Context, imageryRequest core_domain.ImageryRequest,
) (
	core_domain.NewImagery, error,
) {
	requestFilter := core_domain.FilterRequest{
		Id:      imageryRequest.Id,
		User_id: imageryRequest.User_id,
	}
	newImagery, err := h.GetHistory(ctx, requestFilter)
	if err != nil {
		return core_domain.NewImagery{}, err
	}
	if len(newImagery) == 0 {
		return core_domain.NewImagery{}, fmt.Errorf(
			"failed to get imagery with id %d: %w",
			imageryRequest.Id, core_errors.ErrBadRequest,
		)
	}
	return newImagery[0], err
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

	var newUser core_domain.NewImagery
	newUser, err = h.AddNewUserRequest(imageryRequest)
	if err != nil {
		return core_domain.ImageryResponce{}, fmt.Errorf("postgres err: %w", err)
	}

	responce := core_domain.ImageryResponce{
		Image:      output,
		NewImagery: newUser,
	}

	return responce, nil
}
