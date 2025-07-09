package cont

// CalcCheckDigit calculates check digit for owner, equipment category ID and serial number.
// This function was optimized for fun and has a suboptimal reading experience.
func CalcCheckDigit(ownerCode string, equipCatID rune, serialNum int) int {
	var n uint32
	var d uint32 = 1

	for _, c := range ownerCode {
		n += d * charValue(uint32(c))
		d <<= 1
	}

	n += d * charValue(uint32(equipCatID))

	s := uint32(serialNum)
	d = 512
	for d >= 16 {
		n += d * (s % 10)
		d >>= 1
		s /= 10
	}
	return int(n % 11)
}

// charValue returns the index of character plus 10.
// A?BCDEFGHIJK?LMNOPQRSTU?VWXYZ
// A=10, (no 11) B=12, C=13, ... , K=21, (no 22) L=23, ...
func charValue(char uint32) uint32 {
	return char - 55 + (char-56)/10
}
