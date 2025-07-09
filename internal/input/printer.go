package input

// Printer prints inputs and returns nil if no error occurred.
type Printer interface {
	Print(inputs []Input) error
}
