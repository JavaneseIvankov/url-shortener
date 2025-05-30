package app_errors

import (
	"net/http"
)

type AppError struct {
	StatusCode int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:  message,
	}
}

var(
	ErrInternalServerError = NewAppError(
		http.StatusInternalServerError,
		"Something went wrong in our end",
	)
)
