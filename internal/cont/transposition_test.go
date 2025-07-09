package cont

import (
	"fmt"
	"slices"
	"testing"
)

func TestCheckTransposition(t *testing.T) {
	tests := []struct {
		ownerCode  string
		equipCatID rune
		serialNum  int
		checkDigit int
		want       []TpNumber
	}{
		{
			ownerCode:  "ABC",
			equipCatID: 'U',
			serialNum:  123123,
			checkDigit: 7,
			want:       nil,
		},
		{
			ownerCode:  "CMA",
			equipCatID: 'U',
			serialNum:  163912,
			checkDigit: 10,
			want: []TpNumber{
				{Number{"CMA", 'U', 169312, 0}, 2},
				{Number{"CMA", 'U', 163192, 0}, 3},
			},
		},
		{
			ownerCode:  "RCB",
			equipCatID: 'U',
			serialNum:  1130,
			checkDigit: 0,
			want: []TpNumber{
				{Number{"RCB", 'U', 10130, 0}, 1},
			},
		},
		{
			ownerCode:  "WSL",
			equipCatID: 'U',
			serialNum:  801743,
			checkDigit: 10,
			want: []TpNumber{
				{Number{"WSL", 'U', 810743, 0}, 1},
				{Number{"WSL", 'U', 807143, 0}, 2},
				{Number{"WSL", 'U', 801740, 3}, 5},
			},
		},
		{
			ownerCode:  "APL",
			equipCatID: 'U',
			serialNum:  689473,
			checkDigit: 10,
			want: []TpNumber{
				{Number{"APL", 'U', 869473, 0}, 0},
				{Number{"APL", 'U', 698473, 0}, 1},
				{Number{"APL", 'U', 684973, 0}, 2},
				{Number{"APL", 'U', 689743, 0}, 3},
				{Number{"APL", 'U', 689437, 0}, 4},
				{Number{"APL", 'U', 689470, 3}, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%s %s %06d %d", tt.ownerCode, string(tt.equipCatID), tt.serialNum, tt.checkDigit),
			func(t *testing.T) {
				if got := CheckTransposition(tt.ownerCode, tt.equipCatID, tt.serialNum, tt.checkDigit); !slices.Equal(got, tt.want) {
					t.Errorf("CheckTransposition() = %v, want %v", got, tt.want)
				}
			})
	}
}

func BenchmarkCalcCheckTransposition(b *testing.B) {
	for b.Loop() {
		CheckTransposition("APL", 'U', 689473, 10)
	}
}
