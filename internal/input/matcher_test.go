package input

import (
	"testing"
)

func TestMatcher_Match(t *testing.T) {
	match1 := func() Input {
		return Input{
			runeCount: 1,
			matchIndex: func(_ string) []int {
				return []int{0, 1}
			},
		}
	}
	match2 := func() Input {
		return Input{
			runeCount: 2,
			matchIndex: func(_ string) []int {
				return []int{0, 2}
			},
		}
	}
	noMatch := func() Input {
		return Input{
			matchIndex: func(_ string) []int {
				return nil
			},
		}
	}

	tests := []struct {
		name          string
		inputPatterns [][]func() Input
		in            string
		wantedLen     int
	}{
		{
			"Use first pattern",
			[][]func() Input{
				{match1},
				{match1, match2},
			},
			"a",
			1,
		},
		{
			"Use first pattern as default",
			[][]func() Input{
				{noMatch},
				{match2, noMatch},
			},
			"abcd",
			1,
		},
		{
			"Use first best match",
			[][]func() Input{
				{noMatch},
				{match1, noMatch},
				{match1, match1, match1},
			},
			"abcd",
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if patternIdx := Match(tt.in, tt.inputPatterns); len(patternIdx) != tt.wantedLen {
				t.Errorf("Match() = %v, want length %v", patternIdx, tt.wantedLen)
			}
		})
	}
}
