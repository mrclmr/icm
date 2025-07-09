package cont

import "fmt"

// Owner has a code and associated company with its location in the form of country and city.
type Owner struct {
	Code    string
	Company string
	City    string
	Country string
}

// IsOwnerCode checks if string is three upper case letters.
func IsOwnerCode(code string) error {
	if len(code) != 3 {
		return NewValidateError(fmt.Sprintf("%s is not 3 letters long", code))
	}
	if !isUpperLetter(code) {
		return NewValidateError(fmt.Sprintf("%s is not 3 upper case letters", code))
	}
	return nil
}
