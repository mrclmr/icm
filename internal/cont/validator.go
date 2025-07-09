package cont

import "fmt"

// ValidateError is an error for validation of a container number part.
type ValidateError struct {
	message string
}

// NewValidateError returns a new ValidateError.
func NewValidateError(message string) error {
	return &ValidateError{
		message: message,
	}
}

func (e *ValidateError) Error() string {
	return e.message
}

func isOneUpperAlphanumericChar(code string) error {
	if len(code) != 1 {
		return NewValidateError(fmt.Sprintf("%s is not 1 digit", code))
	}
	if !isUpperAlphanumeric(code) {
		return NewValidateError(
			fmt.Sprintf("%s is not 1 upper case alphanumeric character", code))
	}
	return nil
}

func isUpperAlphanumeric(s string) bool {
	for _, r := range s {
		if !isUpperLetter(string(r)) && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func isUpperLetter(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
