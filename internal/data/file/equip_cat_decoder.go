package file

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mrclmr/icm/internal/cont"
)

const equipCatIDsFileName = "equipment-category-id.json"

//go:embed equipment-category-id.json
var equipCatIDsJSON []byte

// NewEquipCatDecoder writes equipment category ID file to path if it not exists and
// returns a struct that uses this file as a data source.
func NewEquipCatDecoder(path string) (*EquipCatDecoder, error) {
	equipCat := &EquipCatDecoder{}
	pathToEquipCat := filepath.Join(path, equipCatIDsFileName)
	if err := initFile(pathToEquipCat, equipCatIDsJSON); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(pathToEquipCat)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &equipCat.categories); err != nil {
		return nil, err
	}
	for ID := range equipCat.categories {
		if err := cont.IsEquipCatID(ID); err != nil {
			return nil, err
		}
	}
	return equipCat, err
}

type EquipCatDecoder struct {
	categories map[string]string
}

// Decode decodes ID to equipment category ID.
func (ecd *EquipCatDecoder) Decode(ID string) (bool, cont.EquipCat) {
	if val, ok := ecd.categories[ID]; ok {
		return true, cont.NewEquipCatID(ID, val)
	}
	return false, cont.EquipCat{}
}

// AllCatIDs returns all equipment category IDs.
func (ecd *EquipCatDecoder) AllCatIDs() []string {
	keys := make([]string, 0, len(ecd.categories))
	for k := range ecd.categories {
		keys = append(keys, k)
	}
	return keys
}
