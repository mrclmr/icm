package cmd

import (
	"github.com/mrclmr/icm/internal/cont"
)

type dummyOwnerDecodeUpdater struct {
	dummyOwnerDecoder
	dummyOwnerUpdater
}

type dummyOwnerDecoder struct{}

func (dummyOwnerDecoder) Decode(code string) (bool, cont.Owner) {
	if code != "ABC" {
		return false, cont.Owner{}
	}
	return true, cont.Owner{
		Code:    "ABC",
		Company: "some-company",
		City:    "some-city",
		Country: "some-country",
	}
}

type dummyOwnerUpdater struct{}

func (dummyOwnerUpdater) GetAllOwnerCodes() []string {
	return []string{"NAR", "RAN"}
}

func (dummyOwnerUpdater) Update([]cont.Owner) error {
	panic("implement me")
}

type dummyEquipCatDecoder struct{}

func (dummyEquipCatDecoder) Decode(ID string) (bool, cont.EquipCat) {
	return true, cont.EquipCat{
		Value: ID,
		Info:  "some-equip-cat-ID",
	}
}

func (dummyEquipCatDecoder) AllCatIDs() []string {
	return []string{"U"}
}

type dummyLengthDecoder struct{}

func (dummyLengthDecoder) Decode(string) (bool, cont.Length) {
	return true, "some-length"
}

type dummyHeightWidthDecoder struct{}

func (dummyHeightWidthDecoder) Decode(string) (bool, cont.Height, cont.Width) {
	return true, "some-height", "some-width"
}

type dummyTypeDecoder struct{}

func (dummyTypeDecoder) Decode(string) (bool, cont.TypeInfo, cont.GroupInfo) {
	return true, "some-type", "some-group"
}
