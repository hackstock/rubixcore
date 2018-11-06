package api

import (
	"net/http"

	"go.uber.org/zap"
)

// Response represent a response sent to
// clients when request is successful
type Response struct {
	Data interface{} `json:"data"`
	Info string      `json:"info"`
}

func handleError(
	w http.ResponseWriter,
	msg string,
	err error,
	logger *zap.Logger,
	code int,
) {
	logger.Warn(msg, zap.Error(err))
	http.Error(w, err.Error(), code)
}

func handleBadRequest(
	w http.ResponseWriter,
	msg string,
	err error,
	logger *zap.Logger,
) {
	handleError(w, msg, err, logger, http.StatusBadRequest)
}

func handleServerError(
	w http.ResponseWriter,
	msg string,
	err error,
	logger *zap.Logger,
) {
	handleError(w, msg, err, logger, http.StatusInternalServerError)
}
