package cont

import "fmt"

// Width describes width of first code in the specified standard size code.
type Width string

// Height describes height of first code in the specified standard size code.
type Height string

// Length describes length of second code in the specified standard size code.
type Length string

// TypeInfo has information about the specified standard type.
type TypeInfo string

// GroupInfo has information about the specified type group.
type GroupInfo string

// IsLengthCode returns nil if input is one upper case alphanumeric character.
func IsLengthCode(code string) error {
	return isOneUpperAlphanumericChar(code)
}

// IsHeightWidthCode returns nil if input is one upper case alphanumeric character.
func IsHeightWidthCode(code string) error {
	return isOneUpperAlphanumericChar(code)
}

// IsTypeCode returns nil if input is two upper case alphanumeric characters.
func IsTypeCode(code string) error {
	if len(code) != 2 {
		return NewValidateError(fmt.Sprintf("%s is not 2 characters long", code))
	}
	if !isUpperAlphanumeric(code) {
		return NewValidateError(
			fmt.Sprintf("%s is not 2 upper case alphanumeric characters", code))
	}
	return nil
}
