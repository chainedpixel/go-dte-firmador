package errors

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
}

// Error returns the error message
func (e DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Well-known error codes
const (
	CodeCertNotFound        = "801"
	CodeInvalid             = "802"
	CodeNoPublicKey         = "803"
	CodeUncatalogued        = "804"
	CodeRequiredData        = "809"
	CodeJSONToStrConversion = "810"
	CodeStrToJSONConversion = "811"
	CodeFileNotFound        = "812"
	CodePasswordInvalid     = "813"
)

// NewDomainError creates a new domain error with the given message and code
func NewDomainError(msg string, code string) DomainError {
	return DomainError{
		Code:    code,
		Message: msg,
	}
}

// NewRequiredDataError creates a new required data error
func NewRequiredDataError(msg string) DomainError {
	return DomainError{
		Code:    CodeRequiredData,
		Message: msg,
	}
}

// NewPasswordInvalidError creates a new password invalid error
func NewPasswordInvalidError(msg string) DomainError {
	return DomainError{
		Code:    CodePasswordInvalid,
		Message: msg,
	}
}
