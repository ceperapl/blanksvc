package common

import (
	"errors"
	"net/http"

	"github.com/company/blanksvc/pkg/repository/postgres"
	"github.com/company/blanksvc/pkg/service"
)

func ErrorToCode(err error) int {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return http.StatusBadRequest
	case errors.Is(err, postgres.ErrTaskNotFound):
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func ErrorKind(err error) string {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return "ErrEmptyString"
	default:
		return "ErrUndefined"
	}
}
