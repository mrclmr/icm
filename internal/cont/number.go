package cont

// Number is a container number with needed properties to conform to the specified standard.
type Number struct {
	OwnerCode    string
	EquipCatID   rune
	SerialNumber int
	CheckDigit   int
}
