package cont

import (
	"slices"
	"testing"
)

func TestCheckTransposition(t *testing.T) {
	type args struct {
		ownerCode  string
		equipCatID rune
		serialNum  int
		checkDigit int
	}
	tests := []struct {
		name string
		args args
		want []TpNumber
	}{
		{
			name: "Test ABC U 123123 7",
			args: args{
				ownerCode:  "ABC",
				equipCatID: 'U',
				serialNum:  123123,
				checkDigit: 7,
			},
			want: nil,
		},
		{
			name: "Test CMA U 163912 (1)0",
			args: args{
				ownerCode:  "CMA",
				equipCatID: 'U',
				serialNum:  163912,
				checkDigit: 10,
			},
			want: []TpNumber{
				{Number{"CMA", 'U', 169312, 0}, 2},
				{Number{"CMA", 'U', 163192, 0}, 3},
			},
		},
		{
			name: "Test RCB U 001130 0",
			args: args{
				ownerCode:  "RCB",
				equipCatID: 'U',
				serialNum:  1130,
				checkDigit: 0,
			},
			want: []TpNumber{
				{Number{"RCB", 'U', 10130, 0}, 1},
			},
		},
		{
			name: "Test WSL U 801743 (1)0",
			args: args{
				ownerCode:  "WSL",
				equipCatID: 'U',
				serialNum:  801743,
				checkDigit: 10,
			},
			want: []TpNumber{
				{Number{"WSL", 'U', 810743, 0}, 1},
				{Number{"WSL", 'U', 807143, 0}, 2},
				{Number{"WSL", 'U', 801740, 3}, 5},
			},
		},
		{
			name: "Test APL U 689473 (1)0",
			args: args{
				ownerCode:  "APL",
				equipCatID: 'U',
				serialNum:  689473,
				checkDigit: 10,
			},
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
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckTransposition(tt.args.ownerCode, tt.args.equipCatID, tt.args.serialNum, tt.args.checkDigit); !slices.Equal(got, tt.want) {
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
