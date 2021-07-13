package domain

import "fmt"

type (
	ErrorCode   string
	ErrorStatus string
)

const (
	ErrSpecValidation     ErrorStatus = "SPEC_VALIDATION_ERROR"
	ErrSpecValidationCode ErrorCode   = "40000"
	ErrDataExists         ErrorStatus = "DATA_EXISTS"
	ErrDataExistsCode     ErrorCode   = "40001"
	ErrInvalidRequest     ErrorStatus = "BAD_REQUEST"
	ErrInvalidRequestCode ErrorCode   = "40002"
	ErrInternalError      ErrorStatus = "INTERNAL_SERVER_ERROR"
	ErrInternalErrorCode  ErrorCode   = "50001"
)

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

func NewError(code ErrorCode, status ErrorStatus, message string) *Error {
	return &Error{
		Code:    string(code),
		Status:  string(status),
		Message: message,
	}
}

type SpecError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
	Errors  []Error
}

func (eg SpecError) Error() string {
	message := eg.Message
	for i := 0; i < len(eg.Errors); i++ {
		message += fmt.Sprintf("\n%s", eg.Errors[i].Error())
	}
	return message
}

type SpecErrorBuilder struct {
	Errors []Error
}

func (eg *SpecErrorBuilder) Append(err *Error) {
	eg.Errors = append(eg.Errors, *err)
}

func (eg *SpecErrorBuilder) AppendBadRequest(message string) {
	err := NewBadRequest(message)
	eg.Errors = append(eg.Errors, *err)
}

func (eg *SpecErrorBuilder) GetIfAny() error {
	if len(eg.Errors) == 0 {
		return nil
	}
	return SpecError{
		Code:    string(ErrSpecValidationCode),
		Status:  string(ErrSpecValidation),
		Message: "Spec Validation Error",
		Errors:  eg.Errors,
	}
}

var NewSpecBuilder = func() *SpecErrorBuilder {
	return &SpecErrorBuilder{
		Errors: make([]Error, 0),
	}
}

var (
	NewErrorDataExists    = func(message string) *Error { return NewError(ErrDataExistsCode, ErrDataExists, message) }
	NewBadRequest         = func(message string) *Error { return NewError(ErrInvalidRequestCode, ErrInvalidRequest, message) }
	NewErrorInternalError = func() *Error { return NewError(ErrInternalErrorCode, ErrInternalError, "") }
)
