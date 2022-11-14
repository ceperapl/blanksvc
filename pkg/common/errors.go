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
	ErrTypeAssertion = errors.New("failed type assertion")
)

func ErrorToCode(err error) int {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func ErrorKind(err error) string {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		return "ErrTaskNotFound"
	default:
		return "ErrUndefined"
	}
}
