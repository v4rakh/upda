package server

import (
	"errors"
	"fmt"
)

var (
	errorValidationNotEmpty            = newServiceError(IllegalArgument, errors.New("assert: empty values are not allowed"))
	errorValidationNotBlank            = newServiceError(IllegalArgument, errors.New("assert: blank values are not allowed"))
	errorValidationPageGreaterZero     = newServiceError(IllegalArgument, errors.New("assert: page has to be greater 0"))
	errorValidationPageSizeGreaterZero = newServiceError(IllegalArgument, errors.New("assert: pageSize has to be greater 0"))

	errorResourceNotFound     = newServiceError(NotFound, errors.New("resource not found"))
	errorResourceAccessDenied = newServiceError(Forbidden, errors.New("resource access denied"))

	errorDatabaseRowsExpected = newServiceDatabaseError(errors.New("action failed, expected affected rows, but got none"))
)

type ErrorCode string

const (
	IllegalArgument ErrorCode = "IllegalArgument"
	Unauthorized    ErrorCode = "Unauthorized"
	Forbidden       ErrorCode = "Forbidden"
	NotFound        ErrorCode = "NotFound"
	Conflict        ErrorCode = "Conflict"
	General         ErrorCode = "General"
)

// newServiceError returns an error that formats as the given text and aligns with builtin error
func newServiceError(status ErrorCode, err error) error {
	return &serviceError{status, fmt.Errorf("service error (%v): %w", status, err)}
}

// newServiceDatabaseError returns an error that formats as the given text and aligns with builtin error
func newServiceDatabaseError(error error) error {
	return newServiceError(General, fmt.Errorf("database error: %w", error))
}

type serviceError struct {
	Status ErrorCode
	Cause  error
}

func (e *serviceError) Error() string {
	return fmt.Sprintf("%v", e.Cause)
}
