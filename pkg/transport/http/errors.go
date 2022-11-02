package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/company/blanksvc/pkg/common"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(common.ErrorToCode(err))
	// nolint: errcheck
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
