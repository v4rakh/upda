package server

import (
	"errors"
	"fmt"
)

var (
	errorValidationNotEmpty              = newServiceError(illegalArgument, errors.New("assert: empty values are not allowed"))
	errorValidationNotBlank              = newServiceError(illegalArgument, errors.New("assert: blank values are not allowed"))
	errorValidationPageGreaterZero       = newServiceError(illegalArgument, errors.New("assert: page has to be greater 0"))
	errorValidationPageSizeGreaterZero   = newServiceError(illegalArgument, errors.New("assert: pageSize has to be greater 0"))
	errorValidationLimitGreaterZero      = newServiceError(illegalArgument, errors.New("assert: limit has to be greater 0"))
	errorValidationSizeGreaterZero       = newServiceError(illegalArgument, errors.New("assert: size has to be greater 0"))
	errorValidationMaxRetriesGreaterZero = newServiceError(illegalArgument, errors.New("assert: max retries has to be greater 0"))

	errorResourceNotFound     = newServiceError(notFound, errors.New("resource not found"))
	errorResourceAccessDenied = newServiceError(forbidden, errors.New("resource access denied"))

	errorDatabaseRowsExpected = newServiceDatabaseError(errors.New("action failed, expected affected rows, but got none"))
)

type errorCode string

const (
	illegalArgument  errorCode = "IllegalArgument"
	unauthorized     errorCode = "Unauthorized"
	forbidden        errorCode = "Forbidden"
	notFound         errorCode = "NotFound"
	methodNotAllowed errorCode = "MethodNotAllowed"
	conflict         errorCode = "Conflict"
	general          errorCode = "General"
)

// newServiceError returns an error that formats as the given text and aligns with builtin error
func newServiceError(status errorCode, err error) error {
	return &serviceError{status, fmt.Errorf("service error (%v): %w", status, err)}
}

// newServiceDatabaseError returns an error that formats as the given text and aligns with builtin error
func newServiceDatabaseError(error error) error {
	return newServiceError(general, fmt.Errorf("database error: %w", error))
}

type serviceError struct {
	Status errorCode
	Cause  error
}

func (e *serviceError) Error() string {
	return fmt.Sprintf("%v", e.Cause)
}
