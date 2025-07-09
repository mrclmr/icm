package input

import (
	"encoding/csv"
)

// Datum represents a datum that is be used by CSVPrinter.
type Datum struct {
	header string
	value  string
}

// NewDatum returns a new Datum.
func NewDatum(header string) Datum {
	return Datum{header: header}
}

// WithValue sets value and returns Datum.
func (d Datum) WithValue(value string) Datum {
	d.value = value
	return d
}

// CSVPrinter prints the set record. Use SetRecord to set a record.
type CSVPrinter struct {
	csvWriter     *csv.Writer
	headers       []string
	record        []string
	headerPrinted bool
	noHeader      bool
}

// NewCSVPrinter creates a new CSVPrinter.
func NewCSVPrinter(csvWriter *csv.Writer, noHeader bool) *CSVPrinter {
	return &CSVPrinter{
		csvWriter: csvWriter,
		noHeader:  noHeader,
	}
}

// Print writes set record to passed writer.
// No header is printed if noHeader is set to false.
// Print returns an error if writing to writer fails.
func (cp *CSVPrinter) Print(inputs []Input) error {
	cp.headers = nil
	cp.record = nil
	for _, input := range inputs {
		for _, datum := range input.data {
			cp.headers = append(cp.headers, datum.header)
			cp.record = append(cp.record, datum.value)
		}
	}

	if !cp.noHeader && !cp.headerPrinted {
		err := cp.csvWriter.Write(cp.headers)
		if err != nil {
			return err
		}
		cp.headerPrinted = true
	}
	err := cp.csvWriter.Write(cp.record)
	if err != nil {
		return err
	}
	cp.csvWriter.Flush()
	return nil
}
