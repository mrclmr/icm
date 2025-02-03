package file

import (
	"encoding/csv"
	"io"

	"github.com/mrclmr/icm/cont"
)

// WriteOwnersCSV accepts a slice of owners and writes CSV to out.
func WriteOwnersCSV(newOwners []cont.Owner, out io.Writer) error {
	csvWriter := csv.NewWriter(out)
	csvWriter.Comma = csvSep

	for _, o := range newOwners {
		csvErr := csvWriter.Write([]string{o.Code, o.Company, o.City, o.Country})
		if csvErr != nil {
			return csvErr
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
