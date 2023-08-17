package validate

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrInvalidID occurs when an ID is not in a valid form.
var ErrInvalidID = errors.New("ID is not in its proper form")

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Err    string      `json:"error"`
	Fields FieldErrors `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request through the
// application with web specific context.
type RequestError struct {
	Err    error
	Status int
	Fields FieldErrors
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewInternalServerError(err string) error {
	return &ErrorResponse{err, nil}
}

// NewFieldError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter field level errors.
func NewFieldErrors(fields FieldErrors) error {
	return &RequestError{errors.New("field validation error"), http.StatusBadRequest, fields}
}

// NewFieldError creates a error of type FieldError.
func NewFieldError(field, err string) FieldError {
	return FieldError{field, err}
}

// RequestError implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *RequestError) Error() string {
	return err.Err.Error()
}

// ErrorResponse implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *ErrorResponse) Error() string {
	return err.Err
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// Error implments the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}
