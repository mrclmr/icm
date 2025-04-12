package cont

import (
	"fmt"
	"testing"
)

func TestCalcCheckDigit(t *testing.T) {
	tests := []struct {
		ownerCode  string
		equipCatID rune
		serialNum  int
		want       int
	}{
		{
			"ABC", 'U', 123456,
			0,
		},
		{
			"NYK", 'U', 8685,
			2,
		},
		{
			"NYK", 'U', 0,
			10,
		},
		{
			"CMA", 'U', 163912,
			10,
		},
		{
			"CMA", 'U', 169312,
			0,
		},
		{
			"CSQ", 'U', 305438,
			3,
		},
		{
			"CSQ", 'U', 999998,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%s %s %06d", tt.ownerCode, string(tt.equipCatID), tt.serialNum),
			func(t *testing.T) {
				if got := CalcCheckDigit(tt.ownerCode, tt.equipCatID, tt.serialNum); got != tt.want {
					t.Errorf("CalcCheckDigit() = %v, want %v", got, tt.want)
				}
			})
	}
}

func BenchmarkCalcCheckDigit(b *testing.B) {
	for b.Loop() {
		CalcCheckDigit("CSQ", 'U', 305438)
	}
}
