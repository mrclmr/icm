package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	// Needed for package embed.
	_ "embed"

	"github.com/mrclmr/icm/cont"
)

const typeFileName = "type.json"

//go:embed type.json
var typeJSON []byte

const groupFileName = "group.json"

//go:embed group.json
var groupJSON []byte

// TypeAndGroupDecoder holds types and groups for decoding.
type TypeAndGroupDecoder struct {
	types  map[string]string
	groups map[string]string
}

// NewTypeDecoder writes type and group file to path if it not exists and
// returns a struct that uses this file as a data source.
func NewTypeDecoder(path string) (*TypeAndGroupDecoder, error) {
	pathToType := filepath.Join(path, typeFileName)
	if err := initFile(pathToType, typeJSON); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(pathToType)
	if err != nil {
		return nil, err
	}

	typeAndGroup := &TypeAndGroupDecoder{}
	if err := json.Unmarshal(b, &typeAndGroup.types); err != nil {
		return nil, err
	}

	pathToGroup := filepath.Join(path, groupFileName)
	if err := initFile(pathToGroup, groupJSON); err != nil {
		return nil, err
	}
	b, err = os.ReadFile(pathToGroup)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &typeAndGroup.groups); err != nil {
		return nil, err
	}
	for typeCode := range typeAndGroup.types {
		if err := cont.IsTypeCode(typeCode); err != nil {
			return nil, err
		}
	}
	return typeAndGroup, nil
}

// Decode returns type and group information for the type code.
func (tgd *TypeAndGroupDecoder) Decode(code string) (bool, cont.TypeInfo, cont.GroupInfo) {
	typeInfoStr, typeFound := tgd.types[code]
	typeInfo := cont.TypeInfo(typeInfoStr)

	if !typeFound {
		return false, "", ""
	}

	groupInfoStr, groupFound := tgd.groups[string(code[0])]
	groupInfo := cont.GroupInfo(groupInfoStr)

	if !groupFound {
		return false, "", ""
	}

	return true, typeInfo, groupInfo
}
