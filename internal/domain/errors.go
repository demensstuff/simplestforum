package domain

import (
	"errors"
	"fmt"
	"net"

	"github.com/gocraft/dbr"
	"github.com/lib/pq"
)

// ErrCode represents a generic error code.
type ErrCode uint8

const (
	ErrCodeInternal           ErrCode = iota + 1 // 1
	ErrCodeAlreadyExists                         // 2
	ErrCodeNotFound                              // 3
	ErrCodeValidation                            // 4
	ErrCodeDatabaseFailure                       // 5
	ErrCodeDatabaseError                         // 6
	ErrCodeNotAuthorized                         // 7
	ErrCodeForbidden                             // 8
	ErrCodeAuthorized                            // 9
	ErrCodeInvalidCredentials                    // 10
	ErrCodeRestricted                            // 11
)

var (
	ErrInternal           = &Error{Code: ErrCodeInternal}
	ErrAlreadyExists      = &Error{Code: ErrCodeAlreadyExists}
	ErrNotFound           = &Error{Code: ErrCodeNotFound}
	ErrValidation         = &Error{Code: ErrCodeValidation}
	ErrDatabaseFailure    = &Error{Code: ErrCodeDatabaseFailure}
	ErrDatabaseError      = &Error{Code: ErrCodeDatabaseError}
	ErrNotAuthorized      = &Error{Code: ErrCodeNotAuthorized, ErrorMessage: "Unauthorized"}
	ErrForbidden          = &Error{Code: ErrCodeForbidden, ErrorMessage: "You don't have permissions to do it"}
	ErrAuthorized         = &Error{Code: ErrCodeAuthorized, ErrorMessage: "You're already logged in"}
	ErrInvalidCredentials = &Error{Code: ErrCodeInvalidCredentials}
	ErrRestricted         = &Error{Code: ErrCodeRestricted, ErrorMessage: "You are restricted from doing it"}
)

// Error stores the information about an error.
type Error struct {
	UUID         string  `json:"uuid"`
	UserID       int64   `json:"user_id"`
	Code         ErrCode `json:"code"`
	ErrorMessage string  `json:"error_message"`
	parent       error
}

// Error returns the error message.
func (e Error) Error() string {
	return e.ErrorMessage
}

// Is returns true if the Error contains the target error.
func (e Error) Is(target error) bool {
	var err *Error

	if !errors.As(target, &err) {
		return false
	}

	return e.Code == err.Code
}

// SetErrorMessage sets a custom error message.
func (e *Error) SetErrorMessage(format string, params ...interface{}) {
	e.ErrorMessage = fmt.Sprintf(format, params...)
}

// NewErrorWrap wraps an external error.
func NewErrorWrap(err error, code ErrCode, errorMessageFormat string, params ...interface{}) *Error {
	e := NewError(code, errorMessageFormat, params...)
	e.parent = err

	return e
}

// NewError creates a new custom Error.
func NewError(code ErrCode, errorMessageFormat string, params ...interface{}) *Error {
	err := &Error{
		Code: code,
	}
	err.SetErrorMessage(errorMessageFormat, params...)

	return err
}

// NewDBErrorWrap wraps the most common database produced errors.
func NewDBErrorWrap(err error) error {
	var opError *net.OpError

	if errors.As(err, &opError) {
		return NewErrorWrap(err, ErrCodeDatabaseFailure, "Database failure: %v", err)
	}

	var pqError *pq.Error

	if errors.As(err, &pqError) {
		return NewErrorWrap(err, ErrCodeDatabaseError, "Database error: %v", err)
	}

	if errors.Is(err, dbr.ErrNotFound) {
		return NewErrorWrap(err, ErrCodeNotFound, "Not found")
	}

	return err
}
