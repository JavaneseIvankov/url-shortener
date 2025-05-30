package pkg

import (
	"encoding/json"
	"errors"
	"javaneseivankov/url-shortener/internal/app_errors"
	"log"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, app_errors.ErrInternalServerError.Error(), app_errors.ErrInternalServerError.StatusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
   if _, writeErr := w.Write(js); writeErr != nil {
		// FIXME: gak yakin
		log.Printf("failed to write response: %v", err)
    }
}

func SendError(w http.ResponseWriter, err error) {
	SendResponse(w, nil, 500, err)	
}

func SendResponse(w http.ResponseWriter, payload interface{}, statusCode int, err error) {
	if err != nil || payload == nil {
		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			payload = map[string]any {
				"error": appErr,
			}
		} else {
			statusCode = app_errors.ErrInternalServerError.StatusCode
			payload = map[string]any {
				"error": app_errors.ErrInternalServerError,
			}
		}
	}

	SendJSON(w, statusCode, payload)
}