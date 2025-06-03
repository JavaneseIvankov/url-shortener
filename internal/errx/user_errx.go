package errx

import "net/http"

var (
	ErrEmailDoesntExist = NewAppError(
		http.StatusInternalServerError,
		"Couldn't find email, email doesn't exist",
	)
	ErrUserIdDoesntExist = NewAppError(
		http.StatusInternalServerError,
		"Couldn't find userID, userID doesn't exist",
	)
)