package exception

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var (
	// Custom error codes
	ErrCodeParseTime       = "PARSE_TIME_ERROR"
	ErrCodeUniqueViolation = "23505"
	ErrNotFound            = pgx.ErrNoRows
	ErrCodeDatabaseError   = "DATABASE_ERROR"
	ErrCodeValidation      = "VALIDATION_ERROR"
)

func NewParseTimeError(field string, err error) *AppError {
	return Wrap(err, http.StatusBadRequest, ErrCodeParseTime, fmt.Sprintf("Invalid time format for %s", field))
}

func NewUniqueViolationError(constraint string, err error) *AppError {
	return Wrap(err, http.StatusConflict, ErrCodeUniqueViolation, fmt.Sprintf("Duplicate entry for constraint: %s", constraint))
}

func NewDatabaseError(operation string, err error) *AppError {
	return Wrap(err, http.StatusInternalServerError, ErrCodeDatabaseError, fmt.Sprintf("Database error during %s", operation))
}

func NewValidationError(message string) *AppError {
	return Wrap(nil, http.StatusBadRequest, ErrCodeValidation, message)
}

func NewNotFoundError(message string) *AppError {
	return Wrap(nil, http.StatusNotFound, ErrNotFound.Error(), message)
}
