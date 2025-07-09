package cont

import "fmt"

// EquipCat has an ID and additional information for the ID.
type EquipCat struct {
	Value string
	Info  string
}

// EquipCatID is the equipment category ID in a container number.
type EquipCatID struct {
	Value string
}

// NewEquipCatID creates a new equipment category with an ID and an information.
func NewEquipCatID(id string, info string) EquipCat {
	return EquipCat{id, info}
}

// IsEquipCatID checks if string is one upper case letter.
func IsEquipCatID(ID string) error {
	if len(ID) != 1 {
		return NewValidateError(fmt.Sprintf("%s is not 1 letter long", ID))
	}
	if !isUpperLetter(ID) {
		return NewValidateError(fmt.Sprintf("%s is not 1 upper case letter", ID))
	}
	return nil
}
