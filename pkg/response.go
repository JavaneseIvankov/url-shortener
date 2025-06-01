package pkg

import (
	"encoding/json"
	"errors"
	"javaneseivankov/url-shortener/internal/errx"
	_logger "javaneseivankov/url-shortener/pkg/logger"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		_logger.Error("response.SendJSON: Failed to marshall value", "value", v)
		http.Error(w, errx.ErrInternalServerError.Error(), errx.ErrInternalServerError.StatusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
   if _, writeErr := w.Write(js); writeErr != nil {
		_logger.Error("response.SendJSON: Failed to write response", "error", err)
    }
}

func SendError(w http.ResponseWriter, err error) {
	SendResponse(w, nil, 500, err)	
}

func SendResponse(w http.ResponseWriter, payload interface{}, statusCode int, err error) {
	if err != nil || payload == nil {
		var appErr *errx.AppError
		if errors.As(err, &appErr) {
			payload = map[string]any {
				"error": appErr,
			}
		} else {
			statusCode = errx.ErrInternalServerError.StatusCode
			payload = map[string]any {
				"error": errx.ErrInternalServerError,
			}
		}
	}

	SendJSON(w, statusCode, payload)
}