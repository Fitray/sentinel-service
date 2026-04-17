package core_errors

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrInvalidArg = errors.New("invalid arguement")
	ErrBadGateway = errors.New("bad gateway")
)

func GetStatusCode(err error) int {
	switch err {
	case ErrInvalidArg:
		return http.StatusBadRequest
	case context.DeadlineExceeded:
		return http.StatusInternalServerError
	case ErrBadGateway:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
