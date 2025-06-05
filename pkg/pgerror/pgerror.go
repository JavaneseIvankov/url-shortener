package pgerror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	UNIQUE_VIOLATION_CODE = "23505"
	FK_VIOLATION_CODE = "23503"
)

type PgError struct {
	ErrCode string
	ConstraintName string
	Err error
}

type PgErrHandler struct {
	pgErrors []PgError
}

func NewPgErrHandler() *PgErrHandler {
	return &PgErrHandler{
		pgErrors: make([]PgError, 0),
	}
}

func (h *PgErrHandler) AddPgErr(errCode string, constraintName string, err error) *PgErrHandler {
	h.pgErrors = append(h.pgErrors, PgError{
		ErrCode: errCode,
		ConstraintName: constraintName,
		Err: err,
	})
	return h
}

func (h *PgErrHandler) Handle(err error) error {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)

	if !ok {
		return err;
	}

	for _, e := range h.pgErrors {
		if pgErr.Code == e.ErrCode && pgErr.ConstraintName == e.ConstraintName {
			return e.Err
		}
	}
	
	return err;
}
