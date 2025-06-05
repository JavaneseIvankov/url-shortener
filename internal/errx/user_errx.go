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
	ErrUserEmailAlreadyExists = NewAppError(
		http.StatusInternalServerError,
		"Email already exists, please use another email!",
	)
	// XXX: this is weird as hell 
	ErrUserIdAlreadyExists = NewAppError(
		http.StatusInternalServerError,
		"Id already exists, please use another email!",
	)
)