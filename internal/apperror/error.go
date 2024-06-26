package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ErrorCode string

const (
	InternalError        ErrorCode = "RAT-005000"
	BadError             ErrorCode = "RAT-004000"
	NotFoundError        ErrorCode = "RAT-004000"
	ConflictError        ErrorCode = "RAT-004000"
	BadGatewayError      ErrorCode = "RAT-005000"
	ForbiddenError       ErrorCode = "RAT-004000"
	TooManyRequestsError ErrorCode = "RAT-004000"
)

var (
	ErrNotFound        = NewAppError(nil, "Not Found", "", string(NotFoundError))
	ConflictErr        = NewAppError(nil, "Conflict", "", string(ConflictError))
	BadGatewayErr      = NewAppError(nil, "Bad Gateway", "", string(BadGatewayError))
	ForbiddenErr       = NewAppError(nil, "Forbidden", "", string(ForbiddenError))
	TooManyRequestsErr = NewAppError(nil, "Too Many Requests", "", string(TooManyRequestsError))
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func InternalServerError(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(InternalError))
}

func BadRequest(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(BadError))
}

func TooManyRequests(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(TooManyRequestsError))
}

func Forbidden(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(ForbiddenError))
}

func Conflict(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(ConflictError))
}

func NotFound(message, developerMessage string) *AppError {
	return NewAppError(fmt.Errorf(message), message, developerMessage, string(NotFoundError))
}

func SystemError(err error) *AppError {
	return NewAppError(err, "internal system error", err.Error(), string(InternalError))
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	ok := errors.As(err, &appErr)
	if !ok {
		return false
	}

	if appErr.Code == string(NotFoundError) {
		return true
	}

	return IsNotFound(appErr.Unwrap())
}
