package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/company/blanksvc/pkg/service"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(errorToCode(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func errorToCode(err error) int {
	switch err {
	case service.ErrEmptyString:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
