package core_responce

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
)

type ResponceHandler struct {
	http.ResponseWriter
	StatusCode int
	logger     core_logger.Logger
}

func NewResponceHandler(
	w http.ResponseWriter,
	logger core_logger.Logger,
) *ResponceHandler {
	return &ResponceHandler{
		ResponseWriter: w,
		logger:         logger,
	}
}

func (rw ResponceHandler) JSONResponce(responce map[string]string, statusCode int) {
	rw.WriteHeader(statusCode)
	if errJson := json.NewEncoder(rw.ResponseWriter).Encode(&responce); errJson != nil {
		rw.logger.Error(fmt.Errorf("failed to respond json: %w", errJson))
	}
}

func (rw ResponceHandler) PanicResponce(msg string, p any) {
	err := fmt.Errorf("a panic occurred during app runtime: %v", p)

	responce := map[string]string{
		"Error":   err.Error(),
		"Message": msg,
	}

	rw.JSONResponce(responce, http.StatusInternalServerError)
}

func (rw ResponceHandler) ErrorResponce(err error, statusCode int) {
	responce := map[string]string{
		"Error": err.Error(),
	}

	rw.JSONResponce(responce, statusCode)
	rw.logger.Error(err)
}

func (rw *ResponceHandler) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}
