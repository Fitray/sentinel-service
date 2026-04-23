package core_errors

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrInvalidArg   = errors.New("invalid arguement")
	ErrBadGateway   = errors.New("bad gateway")
	ErrBadRequest   = errors.New("invalid arguement")
	ErrUnauthorized = errors.New("failed to find user")
)

func GetStatusCode(err error) int {
	switch err {
	case ErrInvalidArg:
		return http.StatusBadRequest
	case context.DeadlineExceeded:
		return http.StatusInternalServerError
	case ErrBadGateway:
		return http.StatusBadGateway
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
