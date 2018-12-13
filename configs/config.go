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

package configs

// Name of the config files and keys for configuration and flags.
const (
	Name           = "config"
	NameWithYmlExt = Name + ".yml"
	SepOE          = "sep-owner-equip"
	SepES          = "sep-equip-serial"
	SepSC          = "sep-serial-check"
	SepCS          = "sep-check-size"
	SepST          = "sep-size-type"
	Pattern        = "pattern"
)

// Cfg returns default config.
func Cfg() []byte {
	return []byte(`#  Separators
#
#  ABC U 123456 0   20 G1
#     ↑ ↑      ↑  ↑   ↑
#     │ │      │  │   └─ ` + SepST + `
#     │ │      │  │
#     │ │      │  └─ ` + SepCS + `
#     │ │      │
#     │ │      └─ ` + SepSC + `
#     │ │
#     │ └─ ` + SepES + `
#     │
#     └─ ` + SepOE + `
#
` + SepOE + `: ' '
` + SepES + `: ' '
` + SepSC + `: ' '
` + SepCS + `: '   '
` + SepST + `: ' '

# Pattern matching mode
#                     auto = matches automatically a pattern
#         container-number = matches a container number
#                    owner = matches a three letter owner code
# owner-equipment-category = matches a three letter owner code with equipment category ID
#                size-type = matches length, width+height and type code
` + Pattern + `: auto
`)
}
