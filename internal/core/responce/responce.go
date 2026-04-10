package core_responce

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponceHandler struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponceHandler(w http.ResponseWriter) *ResponceHandler {
	return &ResponceHandler{
		ResponseWriter: w,
	}
}

func (rw ResponceHandler) PanicResponce(msg string, p any) error {
	err := fmt.Errorf("a panic occurred during app runtime: %v", p)

	responce := map[string]string{
		"Error":   err.Error(),
		"Message": msg,
	}

	rw.ResponseWriter.WriteHeader(http.StatusInternalServerError)

	if errJson := json.NewEncoder(rw.ResponseWriter).Encode(&responce); errJson != nil {
		return fmt.Errorf("%w: failed to respond json: %w", err, errJson)
	}

	return err
}

func (rw *ResponceHandler) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}
