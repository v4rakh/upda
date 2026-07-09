package service_error

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	ErrValidationNotEmpty              = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: empty values are not allowed"))
	ErrValidationNotBlank              = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: blank values are not allowed"))
	ErrValidationPageGreaterZero       = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: page has to be greater 0"))
	ErrValidationPageSizeGreaterZero   = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: pageSize has to be greater 0"))
	ErrValidationLimitGreaterZero      = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: limit has to be greater 0"))
	ErrValidationSizeGreaterZero       = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: size has to be greater 0"))
	ErrValidationMaxRetriesGreaterZero = NewServiceError(ErrCodeIllegalArgument, errors.New("assert: max retries has to be greater 0"))
	ErrResourceNotFound                = NewServiceError(ErrCodeNotFound, errors.New("resource not found"))
	ErrResourceAccessDenied            = NewServiceError(ErrCodeForbidden, errors.New("resource access denied"))
	ErrResourceConflict                = NewServiceError(ErrCodeConflict, errors.New("resource already exists"))
	ErrDatabaseRowsExpected            = NewServiceDatabaseError(errors.New("action failed, expected affected rows, but got none"))
	ErrUnauthorized                    = NewServiceError(ErrCodeUnauthorized, errors.New("unauthorized"))
	ErrStateTransitionNotAllowed       = NewServiceError(ErrCodeIllegalArgument, errors.New("state transition not allowed"))
	ErrStateNotFound                   = NewServiceError(ErrCodeIllegalArgument, errors.New("state not found"))
	ErrStateInUse                      = NewServiceError(ErrCodeConflict, errors.New("state is in use by updates and cannot be deleted"))
	ErrInitialStateRequired            = NewServiceError(ErrCodeIllegalArgument, errors.New("at least one initial state is required"))
)

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	ErrCodeIllegalArgument  ErrorCode = "IllegalArgument"
	ErrCodeUnauthorized     ErrorCode = "Unauthorized"
	ErrCodeForbidden        ErrorCode = "Forbidden"
	ErrCodeNotFound         ErrorCode = "NotFound"
	ErrCodeMethodNotAllowed ErrorCode = "MethodNotAllowed"
	ErrCodeConflict         ErrorCode = "Conflict"
	ErrCodeGeneral          ErrorCode = "General"
)

// NewServiceError returns an error that formats as the given text and aligns with builtin error
func NewServiceError(status ErrorCode, err error) error {
	return &ServiceError{status, fmt.Errorf("service error (%v): %w", status, err)}
}

// NewServiceErrorHttp returns an error from HTTP status codes
func NewServiceErrorHttp(status int, err error) error {
	return NewServiceError(toErrorCode(status), err)
}

// NewServiceDatabaseError returns an error that formats as the given text and aligns with builtin error
func NewServiceDatabaseError(err error) error {
	log.Error().Err(err).Msg("database error")
	return NewServiceError(ErrCodeGeneral, errors.New("a database error occurred"))
}

type ServiceError struct {
	Status ErrorCode
	Cause  error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%v", e.Cause)
}

func toErrorCode(status int) ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return ErrCodeIllegalArgument
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusMethodNotAllowed:
		return ErrCodeMethodNotAllowed
	case http.StatusConflict:
		return ErrCodeConflict
	default:
		return ErrCodeGeneral
	}
}
