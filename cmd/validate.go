// Copyright © 2017 Marcel Meyer meyer@synyx.de
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

package cmd

import (
	"github.com/meyermarcel/iso6346/owner"
	"github.com/meyermarcel/iso6346/parser"
	"github.com/meyermarcel/iso6346/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a container number",
	Long: `Validate a container number.

` + sepHelp,
	Example: `
  iso6346 validate 'ABCU 1234560'

  iso6346 validate --` + sepOE + ` '' --` + sepSC + ` '' 'ABCU 1234560'

  iso6346 validate --only-sizetype '20G0'`,
	Args: cobra.ExactArgs(1),
	Run:  validate,
}

var validateOnlyOwner bool
var validateOnlySizeType bool

func init() {
	validateCmd.Flags().BoolVar(&validateOnlyOwner, "only-owner", false, "validate only owner code")
	validateCmd.Flags().BoolVar(&validateOnlySizeType, "only-sizetype", false, "validate only size and type")

	validateCmd.Flags().String(sepOE, "",
		"ABC(*)U1234560  (*) separator between owner code and equipment category id")
	validateCmd.Flags().String(sepES, "",
		"ABCU(*)1234560  (*) separator between equipment category id and serial number")
	validateCmd.Flags().String(sepSC, "",
		"ABCU123456(*)0  (*) separator between serial number and check digit")

	viper.BindPFlag(sepOE, validateCmd.Flags().Lookup(sepOE))
	viper.BindPFlag(sepES, validateCmd.Flags().Lookup(sepES))
	viper.BindPFlag(sepSC, validateCmd.Flags().Lookup(sepSC))

	RootCmd.AddCommand(validateCmd)
}

func validate(cmd *cobra.Command, args []string) {
	if validateOnlyOwner {
		validateOwner(args[0])
	}

	if validateOnlySizeType {
		validateSizeType(args[0])
	}
	
	validateContNum(args[0])

}

func validateOwner(input string) {
	oce := parser.ParseOwnerCodeOptEquipCat(input)

	oce.OwnerCodeIn.Resolve(owner.Resolver(pathToDB))

	ui.PrintOwnerCode(oce, viper.GetString(sepOE))

	if oce.OwnerCodeIn.IsValidFmt() {
		os.Exit(0)
	}
	os.Exit(1)
}

func validateContNum(input string) {
	num := parser.ParseContNum(input)

	num.OwnerCodeIn.Resolve(owner.Resolver(pathToDB))

	ui.PrintContNum(num, ui.Separators{
		viper.GetString(sepOE),
		viper.GetString(sepES),
		viper.GetString(sepSC),
	})

	if num.CheckDigitIn.IsValidCheckDigit {
		os.Exit(0)
	}
	os.Exit(1)
}

func validateSizeType(input string) {
	st := parser.ParseSizeType(input)

	ui.PrintSizeType(st)

	if st.TypeIn.IsValidFmt() {
		os.Exit(0)
	}
	os.Exit(1)
}
