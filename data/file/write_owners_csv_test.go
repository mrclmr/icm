package file

import (
	"bytes"
	"testing"

	"github.com/mrclmr/icm/cont"
)

func TestWriteOwnersCSV(t *testing.T) {
	tests := []struct {
		name      string
		newOwners []cont.Owner
		wantOut   string
		wantErr   bool
	}{
		{
			"",
			[]cont.Owner{{Code: "ABC", Company: "company", City: "city", Country: "country"}},
			"ABC;company;city;country\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := WriteOwnersCSV(tt.newOwners, out)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteOwnersCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("WriteOwnersCSV() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
