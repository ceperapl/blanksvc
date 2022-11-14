package common

import (
	"errors"
	"net/http"
)

var (
	// serivce errors:

	// repository errors:
	ErrTaskNotFound = errors.New("task is not found")
	// endpoint errors:

	// other:
	ErrTypeAssertion = errors.New("failed type assertion")
	ErrValidation    = errors.New("validation error")
)

func ErrorToCode(err error) int {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func ErrorKind(err error) string {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		return "ErrTaskNotFound"
	case errors.Is(err, ErrValidation):
		return "ErrValidation"
	default:
		return "ErrUndefined"
	}
}
