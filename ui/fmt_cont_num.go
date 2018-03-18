// Copyright © 2018 Marcel Meyer meyermarcel@posteo.de
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
	"github.com/meyermarcel/iso6346/parser"
	"strings"
	"fmt"
	"github.com/meyermarcel/iso6346/equip_cat"
)

func fmtParsedContNum(cn parser.ContNum, seps Separators) string {

	b := strings.Builder{}

	additionalSizeType := cn.LengthIn.IsValidFmt()
	b.WriteString(fmtContNum(cn, seps, additionalSizeType))

	if additionalSizeType {
		b.WriteString(fmtCheckMark(cn.CheckDigitIn.IsValidCheckDigit && cn.TypeAndGroupIn.IsValidFmt()))
	} else {
		b.WriteString(fmtCheckMark(cn.CheckDigitIn.IsValidCheckDigit))
	}

	b.WriteString(fmt.Sprintln())

	var texts []PosTxt

	texts = append(texts, ownerCodeTxt(cn.OwnerCodeIn))

	if !cn.EquipCatIdIn.IsValidFmt() {
		texts = append(texts, NewPosHint(len(indent+seps.OwnerEquip)+3, fmt.Sprintf("%s must be %s", underline("equipment category id"), equipCatIdsAsList())))
	}
	if !cn.SerialNumIn.IsValidFmt() {
		texts = append(texts, NewPosHint(len(indent+seps.OwnerEquip+seps.EquipSerial)+6, fmt.Sprintf("%s must be %s", underline("serial number"), bold("6 numbers"))))
	}

	cdPos := len(indent+seps.OwnerEquip+seps.EquipSerial+seps.SerialCheck) + 10
	if !cn.CheckDigitIn.IsValidCheckDigit {
		if cn.IsCheckDigitCalculable() {
			if cn.CheckDigitIn.IsValidFmt() {
				texts = append(texts, NewPosHint(cdPos, fmt.Sprintf("%s is incorrect (correct: %s)", underline("check digit"),
					green(cn.CheckDigitIn.CalcCheckDigit))))
			} else {
				texts = append(texts, NewPosHint(cdPos, fmt.Sprintf("%s must be a %s (correct: %s)", underline("check digit"), bold("number"),
					green(cn.CheckDigitIn.CalcCheckDigit))))
			}
		} else {
			texts = append(texts, NewPosHint(cdPos, fmt.Sprintf("%s is not calculable", underline("check digit"))))
		}
	}

	if additionalSizeType {
		texts = append(texts, lengthTxt(seps.offsetPosForSizeType(), cn.LengthIn))
		texts = append(texts, heightWidthTxt(seps.offsetPosForSizeType(), cn.HeightWidthIn))
		texts = append(texts, typeAndGroupTxt(seps.offsetPosForSizeType(), cn.TypeAndGroupIn, seps.SizeType))
	}

	b.WriteString(fmtTextsWithArrows(texts...))

	return b.String()
}

func fmtContNum(cn parser.ContNum, seps Separators, additionalSizeType bool) string {

	b := strings.Builder{}

	b.WriteString(indent)
	b.WriteString(fmtOwnerCodeIn(cn.OwnerCodeIn))
	b.WriteString(seps.OwnerEquip)
	b.WriteString(fmtIn(cn.EquipCatIdIn))
	b.WriteString(seps.EquipSerial)
	b.WriteString(fmtIn(cn.SerialNumIn))
	b.WriteString(seps.SerialCheck)

	if !cn.CheckDigitIn.IsValidCheckDigit && cn.CheckDigitIn.IsValidFmt() {
		b.WriteString(fmt.Sprintf("%s", yellow(string(cn.CheckDigitIn.Value()))))
	} else {
		b.WriteString(fmtIn(cn.CheckDigitIn.In))
	}

	if additionalSizeType {
		b.WriteString(seps.CheckSize)
		b.WriteString(fmtIn(cn.LengthIn.In))
		b.WriteString(fmtIn(cn.HeightWidthIn.In))
		b.WriteString(seps.SizeType)
		b.WriteString(fmtIn(cn.TypeAndGroupIn.In))
	}

	return b.String()
}

func equipCatIdsAsList() string {
	ujz := equip_cat.Ids
	return fmt.Sprintf("%s, %s or %s", green(string(ujz[0])), green(string(ujz[1])), green(string(ujz[2])))
}
