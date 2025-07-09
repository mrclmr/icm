package input

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func TestCSVPrinter_Print(t *testing.T) {
	tests := []struct {
		name       string
		noHeader   bool
		inputs     []Input
		wantWriter string
	}{
		{
			name:     "Print CSV with header",
			noHeader: false,
			inputs: []Input{
				{
					data: []Datum{
						{header: "header-1", value: "value-1"},
						{header: "header-2", value: "value-2"},
						{header: "header-3", value: "value-3"},
					},
				},
			},
			wantWriter: `header-1,header-2,header-3
value-1,value-2,value-3
`,
		},
		{
			name:     "Print CSV without header",
			noHeader: true,
			inputs: []Input{
				{
					data: []Datum{
						{header: "header-1", value: "value-1"},
						{header: "header-2", value: "value-2"},
					},
				},
			},
			wantWriter: `value-1,value-2
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}

			csvWriter := csv.NewWriter(writer)
			csvPrinter := NewCSVPrinter(csvWriter, tt.noHeader)
			_ = csvPrinter.Print(tt.inputs)

			csvWriter.Flush()

			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
