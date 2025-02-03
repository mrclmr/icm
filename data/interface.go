package data

import (
	"io"

	"github.com/mrclmr/icm/cont"
)

// OwnerDecoder decodes a code to an owner or generates a random owner.
type OwnerDecoder interface {
	Decode(code string) (bool, cont.Owner)

	GetAllOwnerCodes() []string
}

type WriteOwnersCSVFunc func(newOwners []cont.Owner, out io.Writer) error

// EquipCatDecoder decodes an ID to an equipment category.
type EquipCatDecoder interface {
	Decode(ID string) (bool, cont.EquipCat)

	AllCatIDs() []string
}

// LengthDecoder decodes a code to a length.
type LengthDecoder interface {
	Decode(code string) (bool, cont.Length)
}

// HeightWidthDecoder decodes a code to height and width.
type HeightWidthDecoder interface {
	Decode(code string) (bool, cont.Height, cont.Width)
}

// TypeDecoder decodes a code to type and group information.
type TypeDecoder interface {
	Decode(code string) (bool, cont.TypeInfo, cont.GroupInfo)
}

// TimestampUpdater updates a timestamp with an implemented time.
type TimestampUpdater interface {
	Update() error
}
