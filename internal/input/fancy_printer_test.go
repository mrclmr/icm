// Copyright © 2018 Marcel Meyer meyermarcel@posteo.de
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package input

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestFancyPrinterFactory_Build(t *testing.T) {
	t.Run("Test constructor", func(t *testing.T) {
		writer := bufio.NewWriter(nil)
		inputs := make([]Input, 0)
		separators := []string{"sep1", "sep2"}
		indent := "indent"

		printer := &FancyPrinter{
			writer:     writer,
			inputs:     inputs,
			indent:     indent,
			separators: separators,
		}

		got := NewFancyPrinter(writer, inputs)
		got.SetSeparators(separators...)
		got.SetIndent(indent)

		if !reflect.DeepEqual(got, printer) {
			t.Errorf("FancyPrinterFactory.Build() = %v, want %v", got, printer)
		}
	})
}

func TestFancyPrinter_Print(t *testing.T) {
	type fields struct {
		inputs     []Input
		indent     string
		separators []string
	}

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantWriter string
	}{
		{
			"Print one valid element",
			fields{
				inputs: []Input{
					{
						value: "a",
						infos: []Info{{Text: ""}},
					},
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
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New(""),
						value:     "a",
					},
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
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New(""),
						value:     "a",
					},
					{
						runeCount: 2,
						err:       errors.New(""),
						value:     "bc",
					},
					{
						runeCount: 3,
						err:       errors.New(""),
						value:     "def",
					},
				},
			},
			false,
			`
a bc def  ✘
↑  ↑  ↑
│  │  └─ 
│  │
│  └─ 
│
└─ 

`,
		},
		{
			"Print element with indent",
			fields{
				inputs: []Input{
					{
						runeCount: 0,
						value:     "a",
						infos:     []Info{{Text: ""}},
					},
				},
				indent: "+",
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
			fields{
				inputs: []Input{
					{
						runeCount: 4,
						value:     "abcd",
						infos:     []Info{{Text: ""}},
					},
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
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New(""),
						value:     "a",
					},
					{
						runeCount: 2,
						err:       errors.New(""),
						value:     "bc",
					},
				},
				separators: []string{" * ", " - "},
			},
			false,
			`
a * bc  ✘
↑    ↑
│    └─ 
│
└─ 

`,
		}, {
			"Print multiple elements with no separators",
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New(""),
						value:     "a",
					},
					{
						runeCount: 2,
						err:       errors.New(""),
						value:     "bc",
					},
				},
			},
			false,
			`
a bc  ✘
↑  ↑
│  └─ 
│
└─ 

`,
		},
		{
			"Print one invalid element without value",
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New(""),
						value:     "",
					},
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
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						value:     "a",
						infos:     []Info{{Text: "info text"}},
					},
				},
			},
			false,
			`
a  ✔
↑
└─ info text

`,
		},
		{
			"Print error",
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New("error text"),
						value:     "a",
					},
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
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						err:       errors.New("error line"),
						value:     "",
						infos:     []Info{{Text: "info line"}},
					},
				},
			},
			false,
			`
_  ✘
↑
└─ error line
   info line

`,
		},
		{
			"Print multiple infos with multiples lines",
			fields{
				inputs: []Input{
					{
						runeCount: 1,
						value:     "a",
						infos:     []Info{{Text: "line 1"}, {Text: "line 2"}},
					},
					{
						runeCount: 1,
						value:     "b",
						infos:     []Info{{Text: "line 3"}, {Text: "line 4"}},
					},
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
				inputs: []Input{
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
				separators: []string{"---", "‧‧‧"},
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
				inputs:     tt.fields.inputs,
				indent:     tt.fields.indent,
				separators: tt.fields.separators,
			}
			if err := fp.Print(); (err != nil) != tt.wantErr {
				t.Errorf("FancyPrinter.Print() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("NewFancyPrinterFactory() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
