package input

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestFancyPrinterFactory_Build(t *testing.T) {
	t.Run("Test constructor", func(t *testing.T) {
		writer := bufio.NewWriter(nil)
		separators := []string{"sep1", "sep2"}
		indent := "indent"

		printer := &FancyPrinter{
			writer:     writer,
			indent:     indent,
			separators: separators,
		}

		got := NewFancyPrinter(writer)
		got.SetSeparators(separators...)
		got.SetIndent(indent)

		if !reflect.DeepEqual(got, printer) {
			t.Errorf("FancyPrinterFactory.Build() = %v, want %v", got, printer)
		}
	})
}

func TestFancyPrinter_Print(t *testing.T) {
	type fields struct {
		indent     string
		separators []string
	}

	tests := []struct {
		name       string
		fields     fields
		inputs     []Input
		wantErr    bool
		wantWriter string
	}{
		{
			"Print one valid element",
			fields{},
			[]Input{
				{
					value: "a",
					lines: []string{""},
				},
			},
			false,
			`
a  ✔
↑
└─ 

`,
		},
		{
			"Print one invalid element",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New(""),
					value:     "a",
				},
			},
			false,
			`
a  ✘
↑
└─ 

`,
		},
		{
			"Print multiple invalid elements",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New(""),
					value:     "a",
					lines:     []string{"a text"},
				},
				{
					runeCount: 2,
					err:       errors.New(""),
					value:     "bc",
					lines:     []string{"bc text"},
				},
				{
					runeCount: 3,
					err:       nil,
					value:     "def",
				},
			},
			false,
			`
a bc def  ✘
↑  ↑
│  └─ bc text
│
└─ a text

`,
		},
		{
			"Print element with indent",
			fields{
				indent: "+",
			},
			[]Input{
				{
					runeCount: 0,
					value:     "a",
					lines:     []string{""},
				},
			},
			false,
			`
+a  ✔
 ↑
 └─ 

`,
		},
		{
			"Print 4 character long element",
			fields{},
			[]Input{
				{
					runeCount: 4,
					value:     "abcd",
					lines:     []string{""},
				},
			},
			false,
			`
abcd  ✔
  ↑
  └─ 

`,
		},
		{
			"Print multiple elements with more separators than inputs",
			fields{
				separators: []string{" * ", " - "},
			},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New(""),
					value:     "a",
					lines:     []string{"a text"},
				},
				{
					runeCount: 2,
					err:       errors.New(""),
					value:     "bc",
					lines:     []string{"bc text"},
				},
			},
			false,
			`
a * bc  ✘
↑    ↑
│    └─ bc text
│
└─ a text

`,
		},
		{
			"Print multiple elements with no separators",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New(""),
					value:     "a",
					lines:     []string{"a text"},
				},
				{
					runeCount: 2,
					err:       errors.New(""),
					value:     "bc",
					lines:     []string{"bc text"},
				},
			},
			false,
			`
a bc  ✘
↑  ↑
│  └─ bc text
│
└─ a text

`,
		},
		{
			"Print one invalid element without value",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New(""),
					value:     "",
				},
			},
			false,
			`
_  ✘
↑
└─ 

`,
		},
		{
			"Print info",
			fields{},
			[]Input{
				{
					runeCount: 1,
					value:     "a",
					lines:     []string{"text"},
				},
			},
			false,
			`
a  ✔
↑
└─ text

`,
		},
		{
			"Print error",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New("error text"),
					value:     "a",
				},
			},
			false,
			`
a  ✘
↑
└─ error text

`,
		},
		{
			"Print info with multiples lines",
			fields{},
			[]Input{
				{
					runeCount: 1,
					err:       errors.New("error line"),
					value:     "",
					lines:     []string{"text line"},
				},
			},
			false,
			`
_  ✘
↑
└─ error line
   text line

`,
		},
		{
			"Print multiple lines with multiples lines",
			fields{},
			[]Input{
				{
					runeCount: 1,
					value:     "a",
					lines:     []string{"line 1", "line 2"},
				},
				{
					runeCount: 1,
					value:     "b",
					lines:     []string{"line 3", "line 4"},
				},
			},
			false,
			`
a b  ✔
↑ ↑
│ └─ line 3
│    line 4
│
└─ line 1
   line 2

`,
		},
		{
			"Print separators",
			fields{
				separators: []string{"---", "‧‧‧"},
			},
			[]Input{
				{
					value: "a",
				},
				{
					value: "b",
				},
				{
					value: "c",
				},
			},
			false,
			`
a---b‧‧‧c  ✔

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			fp := &FancyPrinter{
				writer:     writer,
				indent:     tt.fields.indent,
				separators: tt.fields.separators,
			}
			if err := fp.Print(tt.inputs); (err != nil) != tt.wantErr {
				t.Errorf("FancyPrinter.Print() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("FancyPrinter.Print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
