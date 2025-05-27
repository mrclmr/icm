package configs

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Config represents the configuration.
type Config struct {
	Map map[string]string
}

// Overwrite overwrites the configuration with command line flags.
func (c *Config) Overwrite(flagSet *pflag.FlagSet) {
	for k := range map[string]bool{
		FlagNames.Pattern:  true,
		FlagNames.NoHeader: true,
		FlagNames.Output:   true,
		FlagNames.SepOE:    true,
		FlagNames.SepES:    true,
		FlagNames.SepSC:    true,
		FlagNames.SepCS:    true,
		FlagNames.SepST:    true,
	} {

		_, exists := c.Map[k]

		if flagSet.Changed(k) || !exists {
			if k == FlagNames.NoHeader {
				value, _ := flagSet.GetBool(k)
				c.Map[k] = fmt.Sprintf("%t", value)
				continue
			}
			value, _ := flagSet.GetString(k)
			c.Map[k] = value
		}
	}
}

// Pattern returns the pattern config.
func (c *Config) Pattern() string {
	return c.Map[FlagNames.Pattern]
}

// NoHeader returns header config.
func (c *Config) NoHeader() bool {
	value, _ := strconv.ParseBool(c.Map[FlagNames.NoHeader])
	return value
}

// Output returns output config.
func (c *Config) Output() string {
	return c.Map[FlagNames.Output]
}

// SepOE returns the separator between owner and equipment category.
func (c *Config) SepOE() string {
	return c.Map[FlagNames.SepOE]
}

// SepES returns the separator between equipment category and serial number.
func (c *Config) SepES() string {
	return c.Map[FlagNames.SepES]
}

// SepSC returns the separator between serial number and check digit.
func (c *Config) SepSC() string {
	return c.Map[FlagNames.SepSC]
}

// SepCS returns the separator between check digit and size.
func (c *Config) SepCS() string {
	return c.Map[FlagNames.SepCS]
}

// SepST returns the separator between size and type.
func (c *Config) SepST() string {
	return c.Map[FlagNames.SepST]
}

// ReadConfig returns the read config.
func ReadConfig(b []byte) (*Config, error) {
	c := Config{
		make(map[string]string),
	}
	err := yaml.Unmarshal(b, &c.Map)
	if err != nil {
		return nil, err
	}
	_, err = strconv.ParseBool(c.Map[FlagNames.NoHeader])
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Name of the config files and keys for configuration and flags.
const (
	ConfigName           = "config"
	ConfigNameWithYmlExt = ConfigName + ".yml"
)

// Names is the structure for the flag names.
type Names struct {
	Pattern  string
	NoHeader string
	Output   string
	SepOE    string
	SepES    string
	SepSC    string
	SepCS    string
	SepST    string
}

// FlagNames has all the flag names.
var FlagNames = Names{
	Pattern:  "pattern",
	NoHeader: "no-header",
	Output:   "output",
	SepOE:    "sep-owner-equip",
	SepES:    "sep-equip-serial",
	SepSC:    "sep-serial-check",
	SepCS:    "sep-check-size",
	SepST:    "sep-size-type",
}

// Values is the structure for the default flag values.
type Values struct {
	Pattern  string
	NoHeader bool
	Output   string
	SepOE    string
	SepES    string
	SepSC    string
	SepCS    string
	SepST    string
}

// DefaultValues has all the default values.
var DefaultValues = Values{
	Pattern:  "auto",
	NoHeader: false,
	Output:   "auto",
	SepOE:    " ",
	SepES:    " ",
	SepSC:    " ",
	SepCS:    "   ",
	SepST:    " ",
}

// DefaultConfig returns default config.
func DefaultConfig() []byte {
	return []byte(`# Pattern matching mode
#                     auto = matches automatically a pattern
#         container-number = matches a container number
#                    owner = matches a three letter owner code
# owner-equipment-category = matches a three letter owner code with equipment category ID
#                size-type = matches length, width+height and type code
` + FlagNames.Pattern + `: ` + DefaultValues.Pattern + `

# Output mode
#  auto = for a single line 'fancy' and for multiple lines 'csv' output 
#   csv = machine readable CSV output
# fancy = human readable fancy output
` + FlagNames.Output + `: ` + DefaultValues.Output + `

# No header for CSV output
` + FlagNames.NoHeader + `: ` + fmt.Sprintf("%t", DefaultValues.NoHeader) + `

#  Separators
#
#  ABC U 123456 0   20 G1
#     ↑ ↑      ↑  ↑   ↑
#     │ │      │  │   └─ ` + FlagNames.SepST + `
#     │ │      │  │
#     │ │      │  └─ ` + FlagNames.SepCS + `
#     │ │      │
#     │ │      └─ ` + FlagNames.SepSC + `
#     │ │
#     │ └─ ` + FlagNames.SepES + `
#     │
#     └─ ` + FlagNames.SepOE + `
#
` + FlagNames.SepOE + `:  '` + DefaultValues.SepOE + `'
` + FlagNames.SepES + `: '` + DefaultValues.SepES + `'
` + FlagNames.SepSC + `: '` + DefaultValues.SepSC + `'
` + FlagNames.SepCS + `:   '` + DefaultValues.SepCS + `'
` + FlagNames.SepST + `:    '` + DefaultValues.SepST + `'
`)
}
