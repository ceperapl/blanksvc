package common

import (
	"errors"
	"net/http"

	repo "github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/service"
)

func ErrorToCode(err error) int {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return http.StatusBadRequest
	case errors.Is(err, repo.ErrTaskNotFound):
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func ErrorKind(err error) string {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return "ErrEmptyString"
	case errors.Is(err, repo.ErrTaskNotFound):
		return "ErrTaskNotFound"
	default:
		return "ErrUndefined"
	}
}
