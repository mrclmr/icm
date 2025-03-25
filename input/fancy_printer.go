package input

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/logrusorgru/aurora/v4"
	"github.com/mattn/go-isatty"
	"github.com/mrclmr/annot"
)

var au *aurora.Aurora

func init() {
	au = aurora.New(aurora.WithColors(os.Getenv("NO_COLOR") == "" && isatty.IsTerminal(os.Stdout.Fd())))
}

// FancyPrinter prints inputs in a fancy manner. Use NewFancyPrinterFactory to instantiate one.
type FancyPrinter struct {
	writer         io.Writer
	indent         string
	separators     []string
	separatorsFunc func(inputs []Input)
}

// NewFancyPrinter creates a FancyPrinter.
func NewFancyPrinter(writer io.Writer) *FancyPrinter {
	return &FancyPrinter{
		writer: writer,
	}
}

// SetIndent sets the indentation for printing.
func (fp *FancyPrinter) SetIndent(indent string) *FancyPrinter {
	fp.indent = indent
	return fp
}

// SetSeparators sets the separators between inputs. Default separator is ' '.
func (fp *FancyPrinter) SetSeparators(separators ...string) {
	fp.separators = separators
}

// SetSeparatorsFunc sets a function that can set the separators depending on inputs.
func (fp *FancyPrinter) SetSeparatorsFunc(separatorsFunc func(inputs []Input)) {
	fp.separatorsFunc = separatorsFunc
}

// Print writes formatted inputs to writer.
func (fp *FancyPrinter) Print(inputs []Input) error {
	if fp.separatorsFunc != nil {
		fp.separatorsFunc(inputs)
	}

	b := &strings.Builder{}
	_, _ = fmt.Fprintln(b)

	b.WriteString(fp.indent)
	pos := len(fp.indent)

	var annots []*annot.Annot

	valid := true
	for idx, input := range inputs {

		b.WriteString(fmtValue(input))

		sep := " "
		switch {
		case idx == len(inputs)-1:
			sep = ""
		case idx > len(inputs)-1:
			sep = " "
		case idx < len(fp.separators):
			sep = fp.separators[idx]
		}
		b.WriteString(sep)

		if input.err != nil || input.lines != nil {
			a := &annot.Annot{
				Col: pos + input.runeCount/2,
			}
			if input.err != nil && input.err.Error() != "" {
				a.AppendLines(input.err.Error())
			}
			a.AppendLines(input.lines...)
			annots = append(annots, a)
		}
		pos += input.runeCount + utf8.RuneCountInString(sep)

		valid = valid && input.err == nil
	}
	b.WriteString(fmtCheckMark(valid))
	_, _ = fmt.Fprintln(b)
	err := annot.Write(b, annots...)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(b)
	_, _ = io.WriteString(fp.writer, b.String())

	return nil
}

func fmtValue(input Input) string {
	if input.err == nil {
		return fmt.Sprint(au.Green(input.value))
	}
	if input.isValidFmt() {
		return fmt.Sprint(au.Red(input.value))
	}
	return fmt.Sprint(au.Red(strings.Repeat("_", input.runeCount)))
}

func fmtCheckMark(valid bool) string {
	b := strings.Builder{}
	b.WriteString("  ")

	if !valid {
		b.WriteString(fmt.Sprint(au.Red("✘")))
		return b.String()
	}
	b.WriteString(fmt.Sprint(au.Green("✔")))
	return b.String()
}
