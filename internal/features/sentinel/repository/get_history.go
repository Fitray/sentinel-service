package sentinel_repository

import (
	"context"
	"fmt"
	"strconv"

	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_errors "github.com/Fitray/sentinel-service/internal/core/errors"
)

func (r *ImageryRepository) GetHistory(
	ctx context.Context, requestFilter core_domain.FilterRequest,
) ([]core_domain.NewImagery, error) {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		r.Timeout,
	)
	defer cancel()

	query := `
	SELECT id, user_id, city, date_from, date_to, bands, scale, dimensions, created_at, updated_at
	FROM app.requests
	WHERE user_id = $1
	`

	args := []interface{}{requestFilter.User_id}
	argPos := 2

	if requestFilter.City != "" {
		query += fmt.Sprintf(" AND city = $%d", argPos)
		args = append(args, requestFilter.City)
		argPos++
	}

	if requestFilter.Id != "" {
		query += fmt.Sprintf(" AND id = $%d", argPos)
		id, err := strconv.Atoi(requestFilter.Id)
		if err == nil {
			args = append(args, id)
			argPos++
		}
	}

	if !requestFilter.From.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, requestFilter.From)
		argPos++
	}

	if !requestFilter.To.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, requestFilter.To)
		argPos++
	}

	rows, err := r.Pool.Query(ctxTimeout, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find rows: %v: %w",
			core_errors.ErrNotFound,
			err,
		)
	}
	defer rows.Close()

	reqHistory := make([]core_domain.NewImagery, 0)

	for rows.Next() {
		var row core_domain.NewImagery
		if err := rows.Scan(
			&row.Id, &row.User_id, &row.City, &row.From, &row.To, &row.Bands, &row.Scale, &row.Dimensions, &row.Created_at, &row.Updated_at,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v: %w",
				core_errors.ErrInvalidFilter,
				err,
			)
		}
		reqHistory = append(reqHistory, row)
	}

	return reqHistory, nil
}

func (r *ImageryRepository) GetRequestFromID(
	ctx context.Context, user_id string, id int,
) (core_domain.NewImagery, error) {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		r.Timeout,
	)
	defer cancel()

	query := `
	SELECT id, user_id, city, date_from, date_to, bands, scale, dimensions, created_at, updated_at
	FROM app.requests
	WHERE user_id = $1 AND id = $2
	`

	row := r.Pool.QueryRow(ctxTimeout, query, user_id, id)

	var resp core_domain.NewImagery

	if err := row.Scan(
		&resp.Id, &resp.User_id, &resp.City, &resp.From,
		&resp.To, &resp.Bands, &resp.Scale, &resp.Dimensions, &resp.Created_at, &resp.Updated_at,
	); err != nil {
		return core_domain.NewImagery{},
			fmt.Errorf("failed to scan row: %w", err)
	}

	return resp, nil
}
