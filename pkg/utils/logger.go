package utils

import (
	"io"

	"github.com/go-kit/log"
)

// NewLogger creates go-kit logger
func NewLogger(w io.Writer) log.Logger {
	logger := log.NewLogfmtLogger(w)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}
