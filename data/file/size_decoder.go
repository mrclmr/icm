package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	// Needed for package embed.
	_ "embed"

	"github.com/mrclmr/icm/cont"
)

const sizeFileName = "size.json"

//go:embed size.json
var lengthHeightWidthJSON []byte

type size struct {
	Length      map[string]string      `json:"length"`
	HeightWidth map[string]heightWidth `json:"heightWidth"`
}

type heightWidth struct {
	Width  string `json:"height"`
	Height string `json:"width"`
}

// NewSizeDecoder writes last update lengths, height and width file to path if it not exists and
// returns a struct that uses this file as a data source.
func NewSizeDecoder(path string) (*LengthDecoder, *HeightWidthDecoder, error) {
	pathToSizes := filepath.Join(path, sizeFileName)
	if err := initFile(pathToSizes, lengthHeightWidthJSON); err != nil {
		return nil, nil, err
	}
	b, err := os.ReadFile(pathToSizes)
	if err != nil {
		return nil, nil, err
	}

	var s size
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, nil, err
	}
	for lengthCode := range s.Length {
		if err := cont.IsLengthCode(lengthCode); err != nil {
			return nil, nil, err
		}
	}
	for heightWidthCode := range s.HeightWidth {
		if err := cont.IsHeightWidthCode(heightWidthCode); err != nil {
			return nil, nil, err
		}
	}
	return &LengthDecoder{s.Length}, &HeightWidthDecoder{s.HeightWidth}, nil
}

// LengthDecoder holds the lengths for decoding.
type LengthDecoder struct {
	lengths map[string]string
}

// Decode returns length for a given length code.
func (ld *LengthDecoder) Decode(code string) (bool, cont.Length) {
	if val, ok := ld.lengths[code]; ok {
		return true, cont.Length(val)
	}
	return false, ""
}

// HeightWidthDecoder holds height and widths for decoding.
type HeightWidthDecoder struct {
	heightWidths map[string]heightWidth
}

// Decode returns height and width for given height and width code.
func (hwd *HeightWidthDecoder) Decode(code string) (bool, cont.Height, cont.Width) {
	if val, ok := hwd.heightWidths[code]; ok {
		return true, cont.Height(val.Height), cont.Width(val.Width)
	}
	return false, "", ""
}
