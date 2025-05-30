package errx

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
	ErrShortLinkUnauthorizedOperation = NewAppError(
		http.StatusInternalServerError,
		"Unauthorized to perform such operation",
	)
)
