package pgerror

import (
	"errors"
	"javaneseivankov/url-shortener/pkg/logger"

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
	logger.Debug("pgerror.AddPgErr: Added new PgErr", "errcode", errCode, "constraint", constraintName, "error", err)
	return h
}

func (h *PgErrHandler) Handle(err error) error {
    var pgErr *pgconn.PgError
    ok := errors.As(err, &pgErr)

	 if err == nil {
		return nil
	 }

    if !ok {
        logger.Debug("pgerror.Handle: Error is not a PgError, returning original error", "error", err.Error())
        return err
    }

    logger.Debug("pgerror.Handle: Handling PgError", "code", pgErr.Code, "constraint", pgErr.ConstraintName)

    for _, e := range h.pgErrors {
        logger.Debug("pgerror.Handle: Checking PgError against registered errors", "registeredCode", e.ErrCode, "registeredConstraint", e.ConstraintName)
        if pgErr.Code == e.ErrCode && pgErr.ConstraintName == e.ConstraintName {
            logger.Debug("pgerror.Handle: Match found, returning mapped error", "code", e.ErrCode, "constraint", e.ConstraintName)
            return e.Err
        }
    }

    logger.Debug("pgerror.Handle: No match found, returning original error", "code", pgErr.Code, "constraint", pgErr.ConstraintName)
    return err
}