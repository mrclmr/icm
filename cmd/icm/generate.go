package cmd

import (
	"fmt"
	"io"
	"math/rand/v2"
	"path/filepath"
	"strconv"

	"github.com/mrclmr/icm/configs"
	"github.com/mrclmr/icm/cont"
	"github.com/mrclmr/icm/data"

	"github.com/spf13/cobra"
)

type countValue struct {
	value int
}

func (r *countValue) String() string {
	return strconv.Itoa(r.value)
}

func (r *countValue) Set(value string) error {
	count, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if count < 1 {
		return fmt.Errorf("%d is not greater than 0", count)
	}
	r.value = count
	return nil
}

func (*countValue) Type() string {
	return "int"
}

type ownerValue struct {
	value string
}

func (o *ownerValue) String() string {
	return o.value
}

func (o *ownerValue) Set(value string) error {
	if err := cont.IsOwnerCode(value); err != nil {
		return err
	}
	o.value = value
	return nil
}

func (*ownerValue) Type() string {
	return "string"
}

type serialNumValue struct {
	value int
}

func (r *serialNumValue) String() string {
	return strconv.Itoa(r.value)
}

func (r *serialNumValue) Set(value string) error {
	serialNum, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if serialNum < 0 || serialNum > 999999 {
		return fmt.Errorf("%d is not in range from 0 to 999999", serialNum)
	}
	r.value = serialNum
	return nil
}

func (*serialNumValue) Type() string {
	return "int"
}

func newGenerateCmd(writer, writerErr io.Writer, config *configs.Config, ownerDecoder data.OwnerDecoder, r *rand.Rand) *cobra.Command {
	count := countValue{value: 1}
	startValue := serialNumValue{}
	endValue := serialNumValue{}
	ownerValue := ownerValue{}
	var excludeCheckDigit10 bool
	var excludeErrorProneSerialNumbers bool

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate unique container numbers",
		Long: `Generated container numbers are unique. Owners specified in

  ` + filepath.Join("$HOME", appDir, "data", "owner.csv") + `

are used. Owners can be updated by 'icm download-owners --help' command.

Equipment category ID 'U' is used for every generated container number.

For a custom owner code use the --owner-code flag.

For a custom serial number use the --start and --end flags and optionally the --count flag.
Using only the --count flag generates pseudo random serial numbers.

` + sepHelp,
		Example: `icm generate
icm generate --count 10
# Generate container numbers with custom format
icm generate --count 10 --sep-owner-equip '' --sep-serial-check '-'
# Generate container numbers without error-prone serial numbers
icm generate --count 10 --exclude-check-digit-10 --exclude-error-prone-serial-numbers
# Generate container numbers within serial number range
icm generate --start 100500 --count 10
icm generate --start 100500 --end 100600
icm generate --start 100500 --end 100600 --owner ABC
# Generate CSV data set
icm generate --count 1000000 | icm validate`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config.Overwrite(cmd.Flags())

			builder := cont.NewUniqueGeneratorBuilder(r).
				Count(count.value).
				ExcludeCheckDigit10(excludeCheckDigit10).
				ExcludeErrorProneSerialNumbers(excludeErrorProneSerialNumbers)

			if cmd.Flags().Changed("owner") {
				builder.OwnerCodes([]string{ownerValue.value})
			} else {
				builder.OwnerCodes(ownerDecoder.GetAllOwnerCodes())
			}

			if cmd.Flags().Changed("start") {
				builder.Start(startValue.value)
			}

			if cmd.Flags().Changed("end") {
				builder.End(endValue.value)
			}

			generator, err := builder.Build()
			if err != nil {
				return err
			}
			for generator.Generate() {
				cn := generator.ContNum()
				_, err := io.WriteString(writer, fmt.Sprintf("%s%s%s%s%06d%s%d\n",
					cn.OwnerCode, config.SepOE(),
					string(cn.EquipCatID), config.SepES(),
					cn.SerialNumber, config.SepSC(),
					cn.CheckDigit))
				writeErr(writerErr, err)
			}
			return nil
		},
	}

	generateCmd.Flags().SortFlags = false

	generateCmd.Flags().VarP(&count, "count", "c", "count of container numbers")
	generateCmd.Flags().VarP(&startValue, "start", "s", "start of serial number range")
	generateCmd.Flags().VarP(&endValue, "end", "e", "end of serial number range")
	generateCmd.Flags().Var(&ownerValue, "owner", "custom owner code")
	generateCmd.Flags().BoolVar(&excludeCheckDigit10, "exclude-check-digit-10", false, "exclude check digit 10")
	generateCmd.Flags().BoolVar(&excludeErrorProneSerialNumbers, "exclude-transposition-errors", false,
		"exclude possible transposition errors")
	err := generateCmd.Flags().MarkDeprecated("exclude-transposition-errors", "use instead --exclude-error-prone-serial-numbers")
	if err != nil {
		return nil
	}
	generateCmd.Flags().BoolVar(&excludeErrorProneSerialNumbers, "exclude-error-prone-serial-numbers", false,
		"exclude error-prone serial numbers. For example swapping the second 0 and first 1 of RCB U 001130 0 results in container number RCB U 010130 0 with a valid check digit 0")

	generateCmd.Flags().String(configs.FlagNames.SepOE, configs.DefaultValues.SepOE,
		"ABC(x)U1234560  (x) separates owner code and equipment category id")
	generateCmd.Flags().String(configs.FlagNames.SepES, configs.DefaultValues.SepES,
		"ABCU(x)1234560  (x) separates equipment category id and serial number")
	generateCmd.Flags().String(configs.FlagNames.SepSC, configs.DefaultValues.SepSC,
		"ABCU123456(x)0  (x) separates serial number and check digit")

	return generateCmd
}
