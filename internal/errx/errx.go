package errx

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
	ErrNoBearerToken = NewAppError(
		http.StatusUnauthorized,
		"No bearer token found",
	)
	ErrBearerTokenInvalidFormat = NewAppError(
		http.StatusUnauthorized,
		"Invalid bearer token format",
	)
	ErrInvalidBearerToken = NewAppError(
		http.StatusUnauthorized,
		"Invalid bearer token",
	)
	ErrExpiredBearerToken = NewAppError(
		http.StatusUnauthorized,
		"Session expired",
	)
	ErrxMalformedBearerToken = NewAppError(
		http.StatusUnauthorized,
		"Bearer token format is not right!",
	)
)
