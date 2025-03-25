package input

import (
	"errors"
	"testing"
)

func TestInputHasCorrectValue(t *testing.T) {
	match1 := func() Input {
		return Input{
			runeCount: 1,
			matchIndex: func(_ string) []int {
				return []int{0, 1}
			},
			validate: func(_ string, _ []string) ([]string, []Datum, error) {
				return []string{"match 1"}, nil, nil
			},
		}
	}
	match2 := func() Input {
		return Input{
			runeCount: 2,
			toUpper:   true,
			matchIndex: func(_ string) []int {
				return []int{0, 2}
			},
			validate: func(_ string, _ []string) ([]string, []Datum, error) {
				return []string{"match 2"}, nil, nil
			},
		}
	}
	validFmtButInvalidMatch := func() Input {
		return Input{
			runeCount: 1,
			matchIndex: func(_ string) []int {
				return []int{0, 1}
			},
			validate: func(_ string, _ []string) ([]string, []Datum, error) {
				return nil, nil, errors.New("")
			},
		}
	}

	type wantedInput struct {
		value          string
		previousValues []string
		err            bool
		lines          []string
	}

	tests := []struct {
		name         string
		newInputs    []func() Input
		in           string
		wantedInputs []wantedInput
		wantErr      bool
	}{
		{
			"Match single value",
			[]func() Input{match1},
			"a",
			[]wantedInput{
				{
					"a",
					nil,
					false,
					[]string{"match 1"},
				},
			},
			false,
		},
		{
			"Match multiple values",
			[]func() Input{match1, match2},
			"abcd",
			[]wantedInput{
				{
					"a",
					nil,
					false,
					[]string{"match 1"},
				},
				{
					"BC",
					[]string{"a"},
					false,
					[]string{"match 2"},
				},
			},
			false,
		},
		{
			"Match but invalid",
			[]func() Input{validFmtButInvalidMatch},
			"a",
			[]wantedInput{
				{
					"a",
					nil,
					true,
					nil,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs, err := Validate(tt.in, tt.newInputs)
			if len(inputs) != len(tt.wantedInputs) {
				t.Errorf("inputs len %v, want %v", len(inputs), len(tt.wantedInputs))
			}
			if (err == nil) == tt.wantErr {
				t.Errorf("got = %v, wantErr is %v", err, tt.wantErr)
			}
			for i, input := range inputs {
				if input.value != tt.wantedInputs[i].value {
					t.Errorf("value is %v, want %v", input.value, tt.wantedInputs[i].value)
				}
				if len(input.previousValues) != len(tt.wantedInputs[i].previousValues) {
					t.Errorf("previousValues len %v, want %v", len(inputs), len(tt.wantedInputs[i].previousValues))
				}
				if input.previousValues != nil {
					for j, previousValue := range input.previousValues {
						if previousValue != tt.wantedInputs[i].previousValues[j] {
							t.Errorf("previous value is %v, want %v", previousValue, tt.wantedInputs[i].previousValues[j])
						}
					}
				}
				if (input.err != nil) != tt.wantedInputs[i].err {
					t.Errorf("err is %v, want %v", input.err != nil, tt.wantedInputs[i].err)
				}
				if len(input.lines) != len(tt.wantedInputs[i].lines) {
					t.Errorf("input lines len %v, want %v", len(input.lines), len(tt.wantedInputs[i].lines))
				}
				if input.lines != nil {
					for j, line := range input.lines {
						if line != tt.wantedInputs[i].lines[j] {
							t.Errorf("line text is %v, want %v", line, tt.wantedInputs[i].lines[j])
						}
					}
				}
			}
		})
	}
}
