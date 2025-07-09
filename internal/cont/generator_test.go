package cont

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"testing"
)

func TestGeneratorBuilder(t *testing.T) {
	type fields struct {
		codes            []string
		count            int
		rangeStart       int
		rangeEnd         int
		exclCheckDigit10 bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *UniqueGenerator
		wantErr bool
	}{
		{
			"Build unique container number generator with random serial iterator",
			fields{
				[]string{"ABC"},
				2,
				-1,
				-1,
				true,
			},
			&UniqueGenerator{
				codes:    []string{"ABC"},
				lenCodes: 1,
				serialNumIt: &randSerialNumIt{
					randOffset: 1812594575390091523,
				},
				count:            2,
				exclCheckDigit10: true,
			},
			false,
		},
		{
			"Build unique container number generator with sequential serial iterator and only range start",
			fields{
				codes:      []string{"ABC"},
				count:      3,
				rangeStart: 2,
				rangeEnd:   -1,
			},
			&UniqueGenerator{
				codes:       []string{"ABC"},
				lenCodes:    1,
				serialNumIt: newSeqSerialNumIt(2),
				count:       3,
			},
			false,
		},
		{
			"Build unique container number generator with sequential serial iterator and only range end",
			fields{
				codes:      []string{"ABC"},
				count:      4,
				rangeStart: -1,
				rangeEnd:   2,
			},
			&UniqueGenerator{
				codes:       []string{"ABC"},
				lenCodes:    1,
				serialNumIt: newSeqSerialNumIt(-1),
				count:       4,
			},
			false,
		},
		{
			"Build unique container number generator with sequential serial iterator and start and end range",
			fields{
				codes:      []string{"ABC"},
				count:      1,
				rangeStart: 2,
				rangeEnd:   5,
			},
			&UniqueGenerator{
				codes:       []string{"ABC"},
				lenCodes:    1,
				serialNumIt: newSeqSerialNumIt(2),
				count:       4,
			},
			false,
		},
		{
			"Build returns error for no owner codes",
			fields{
				count:      1,
				rangeStart: -1,
				rangeEnd:   -1,
			},
			nil,
			true,
		},
		{
			"Build returns error for count less than 1",
			fields{
				codes:      []string{"ABC"},
				rangeStart: -1,
				rangeEnd:   -1,
			},
			nil,
			true,
		},
		{
			"Build returns error for count less than 1 and positive start",
			fields{
				codes:      []string{"ABC"},
				rangeStart: 1,
				rangeEnd:   -1,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rand.New(rand.NewPCG(1, 0))
			gb := NewUniqueGeneratorBuilder(r).
				OwnerCodes(tt.fields.codes).
				Count(tt.fields.count).
				Start(tt.fields.rangeStart).
				End(tt.fields.rangeEnd).
				ExcludeCheckDigit10(tt.fields.exclCheckDigit10)
			got, err := gb.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratorBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratorBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueGenerator(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 0))

	tests := []struct {
		name             string
		generatorBuilder *GeneratorBuilder
		seqSerialNum     bool
		want             int
	}{
		{
			"Generate 3 unique container numbers with random serial numbers",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Count(3),
			false,
			3,
		},
		{
			"Generate 1 container number",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Count(3).
				Start(1).
				End(1),
			true,
			1,
		},
		{
			"Generate 3 unique container numbers with sequential serial numbers and start 1",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Count(3).
				Start(1),
			true,
			3,
		},
		{
			"Generate 4 unique container numbers with sequential serial numbers and end 2",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Count(4).
				End(2),
			true,
			4,
		},
		{
			"Generate 5 unique container numbers with sequential serial numbers and start 1 and end 5",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Start(1).
				End(5),
			true,
			5,
		},
		{
			"Generate 6 unique container numbers with sequential serial numbers and start 999997 and end 2",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Start(999997).
				End(2),
			true,
			6,
		},
		{
			"Generate 1 unique container numbers with sequential serial numbers and exclude possible transposition errors",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"ABC"}).
				Start(801743).
				Count(1).
				ExcludeErrorProneSerialNumbers(true),
			false,
			1,
		},
		{
			"Generate 2000001 unique container numbers with random serial numbers",
			NewUniqueGeneratorBuilder(r).
				OwnerCodes([]string{"AAA", "AAB", "ABB"}).
				Count(2000001),

			false,
			2000001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := tt.generatorBuilder.Build()
			if err != nil {
				t.Errorf("GeneratorBuilder.Build() returned error, %v", err)
				return
			}

			lastNum := g.serialNumIt.num() - 1
			diff := 0
			contNumbers := map[string]bool{}
			for g.Generate() {
				number := g.ContNum().SerialNumber
				if number < 0 || number > 999999 {
					t.Errorf("UniqueGenerator.Generate() generated a serial number out of range, %v", number)
					return
				}
				diff += ((lastNum + 1) % 1000000) - (number)
				lastNum = number
				cn := g.ContNum()
				contNumbers[fmt.Sprintf("%s%s%06d%d\n",
					cn.OwnerCode,
					string(cn.EquipCatID),
					cn.SerialNumber,
					cn.CheckDigit)] = true
			}

			if tt.seqSerialNum && diff != 0 {
				t.Errorf("UniqueGenerator.Generate() generated not sequential serial numbers, diff is %v", diff)
				return
			}

			if got := len(contNumbers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueGenerator.Generate() generated %v unique container numbers, want %v", got, tt.want)
			}
		})
	}
}
