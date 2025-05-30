package app_errors

import "net/http"

var(
	ErrShortLinkNotFound = NewAppError(
		http.StatusInternalServerError,
		"Shortlink not found!",
		)
	ErrShortLinkAlreadyExists = NewAppError(
		http.StatusInternalServerError,
		"Shortlink already exists!",
		)
)
