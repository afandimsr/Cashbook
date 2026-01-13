package apperror

import "net/http"

type AppError struct {
	Code      int
	ErrorCode string
	Message   string
	Err       error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) WithCode(code string) *AppError {
	e.ErrorCode = code
	return e
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Helpers
func BadRequest(msg string, err error) *AppError {
	return New(http.StatusBadRequest, msg, err)
}

func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err)
}

func Unauthorized(msg string, err error) *AppError {
	return New(http.StatusUnauthorized, msg, err)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "internal server error", err)
}
