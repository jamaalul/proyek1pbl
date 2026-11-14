package helper

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int
	Message string
	Errors  map[string][]string
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

func BadRequest(message string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Status: http.StatusUnauthorized, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Status: http.StatusForbidden, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Status: http.StatusNotFound, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{Status: http.StatusConflict, Message: message}
}

func Validation(errors map[string][]string) *AppError {
	return &AppError{Status: http.StatusUnprocessableEntity, Message: "validasi gagal", Errors: errors}
}

func TooManyRequests(message string) *AppError {
	return &AppError{Status: http.StatusTooManyRequests, Message: message}
}

func Internal(cause error) *AppError {
	return &AppError{Status: http.StatusInternalServerError, Message: "terjadi kesalahan pada server", cause: cause}
}
