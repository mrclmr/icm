package cont

import "math"

// TpNumber represents a transposed serial number.
type TpNumber struct {
	Number
	// Pos is the position of the two transposed numbers starting with 0.
	// For example Pos 1 means that the second and third digits are transposed.
	Pos int
}

// CheckTransposition returns an array of TpNumber's for error-prone serial numbers.
// CheckTransposition returns nil if no error-prone serial number is found.
// Not equal adjacent digits including check digit are transposed and checked.
func CheckTransposition(ownerCode string, equipCatID rune, serialNum int, checkDigit int) []TpNumber {
	checkDigit = checkDigit % 10

	var contNums []TpNumber

	// Only container numbers with check digit 0, 10 or 3 are affected
	if checkDigit != 3 && checkDigit != 0 {
		return contNums
	}

	// 5, 4, 3, 2, 1
	for idxRight := 5; idxRight > 0; idxRight-- {
		swapped, transposedSerialNum := swapDigits(serialNum, idxRight-1, idxRight)
		if !swapped {
			continue
		}
		calcCheckDigit := CalcCheckDigit(ownerCode, equipCatID, transposedSerialNum) % 10
		if checkDigit == calcCheckDigit {
			// 0, 1, 2, 3, 4
			pos := idxRight*-1 + 5
			contNums = append(contNums, TpNumber{Number{ownerCode, equipCatID, transposedSerialNum, calcCheckDigit}, pos})
		}
	}

	serialNumLastDigit := serialNum % 10
	if checkDigit == serialNumLastDigit {
		return contNums
	}

	transposedCheckDigitSerialNum := ((serialNum / 10) * 10) + checkDigit
	calcCheckDigit := CalcCheckDigit(ownerCode, equipCatID, transposedCheckDigitSerialNum) % 10
	if serialNumLastDigit == calcCheckDigit {
		contNums = append(contNums, TpNumber{Number{ownerCode, equipCatID, transposedCheckDigitSerialNum, serialNumLastDigit}, 5})
	}
	return contNums
}

// swapDigits returns true if the digits are different and returns the number with swapped numbers.
// false is returned if the digits are same and 0 is returned.
// Position 0 is the digit beginning on the right side of number.
func swapDigits(number int, pos1, pos2 int) (bool, int) {
	p1 := int(math.Pow10(pos1))
	p2 := int(math.Pow10(pos2))

	digit1 := (number / p1) % 10
	digit2 := (number / p2) % 10
	if digit1 == digit2 {
		return false, 0
	}

	number -= digit1 * p1
	number -= digit2 * p2

	number += digit1 * p2
	number += digit2 * p1

	return true, number
}
