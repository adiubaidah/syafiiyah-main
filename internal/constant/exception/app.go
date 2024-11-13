package exception

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type AppError struct {
	Code        int
	ErrCode     string
	Message     string
	OriginalErr error
}

func (e *AppError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("[%s] %s: %v", e.ErrCode, e.Message, e.OriginalErr)
	}
	return fmt.Sprintf("[%s] %s", e.ErrCode, e.Message)
}

func Wrap(err error, code int, errCode string, message string) *AppError {
	return &AppError{
		Code:        code,
		ErrCode:     errCode,
		Message:     message,
		OriginalErr: err,
	}
}

func DatabaseErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
