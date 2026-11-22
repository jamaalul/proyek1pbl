package helper

import (
	"fmt"
	"net/http"

	"github.com/jamaalul/proyek1pbl/app/model"
)

type AppError struct {
	Status  int
	Message string
	Errors  []model.ErrorItem
	cause   error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func NewAppError(status int, code, message string, field ...string) *AppError {
	f := ""
	if len(field) > 0 {
		f = field[0]
	}
	return &AppError{
		Status:  status,
		Message: message,
		Errors: []model.ErrorItem{
			{Field: f, Code: code, Message: message},
		},
	}
}

// PRD §8.2 Constructors with codes
func ErrBadRequest(code, message string, field ...string) *AppError {
	return NewAppError(http.StatusBadRequest, code, message, field...)
}

func ErrUnauthorized(code, message string) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, "FORBIDDEN", message)
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, "NOT_FOUND", message)
}

func ErrConflict(code, message string) *AppError {
	return NewAppError(http.StatusConflict, code, message)
}

func ErrUnprocessable(code, message string, field ...string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, code, message, field...)
}

func ErrInternal(cause error) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Message: "terjadi kesalahan pada server",
		Errors: []model.ErrorItem{
			{Code: "INTERNAL", Message: "terjadi kesalahan pada server"},
		},
		cause: cause,
	}
}

// Backward-compatible constructors
func BadRequest(message string) *AppError {
	return ErrBadRequest("BAD_REQUEST", message)
}

func Unauthorized(message string) *AppError {
	return ErrUnauthorized("UNAUTHORIZED", message)
}

func Forbidden(message string) *AppError {
	return ErrForbidden(message)
}

func NotFound(message string) *AppError {
	return ErrNotFound(message)
}

func Conflict(message string) *AppError {
	return ErrConflict("CONFLICT", message)
}

func TooManyRequests(message string) *AppError {
	return NewAppError(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

func Internal(cause error) *AppError {
	return ErrInternal(cause)
}

func Validation(errorsMap map[string][]string) *AppError {
	var items []model.ErrorItem
	for field, msgs := range errorsMap {
		for _, msg := range msgs {
			items = append(items, model.ErrorItem{
				Field:   field,
				Code:    "VALIDATION_ERROR",
				Message: msg,
			})
		}
	}
	return &AppError{
		Status:  http.StatusUnprocessableEntity,
		Message: "validasi gagal",
		Errors:  items,
	}
}
