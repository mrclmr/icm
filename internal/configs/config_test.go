package configs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		b       []byte
		want    *Config
		wantErr bool
	}{
		{
			"parse default config",
			DefaultConfig(),
			&Config{Map: map[string]string{
				FlagNames.Pattern:  DefaultValues.Pattern,
				FlagNames.NoHeader: fmt.Sprintf("%t", DefaultValues.NoHeader),
				FlagNames.Output:   DefaultValues.Output,
				FlagNames.SepOE:    DefaultValues.SepOE,
				FlagNames.SepES:    DefaultValues.SepES,
				FlagNames.SepSC:    DefaultValues.SepSC,
				FlagNames.SepCS:    DefaultValues.SepCS,
				FlagNames.SepST:    DefaultValues.SepST,
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_OverwriteString(t *testing.T) {
	c := &Config{
		Map: map[string]string{
			FlagNames.Pattern: DefaultValues.Pattern,
		},
	}
	set := &pflag.FlagSet{}
	want := "not default"
	set.String(FlagNames.Pattern, DefaultValues.Pattern, "")
	_ = set.Set(FlagNames.Pattern, want)

	set.String(FlagNames.Output, DefaultValues.Output, "")

	c.Overwrite(set)
	if c.Pattern() != want {
		t.Errorf("Pattern() got = %v, want %v", c.Pattern(), want)
	}
	if c.Output() != DefaultValues.Output {
		t.Errorf("Output() got = %v, want %v", c.Output(), DefaultValues.Output)
	}
}

func TestConfig_OverwriteBool(t *testing.T) {
	c := &Config{
		Map: map[string]string{
			FlagNames.NoHeader: fmt.Sprintf("%t", DefaultValues.NoHeader),
		},
	}
	set := &pflag.FlagSet{}
	set.Bool(FlagNames.NoHeader, DefaultValues.NoHeader, "")
	_ = set.Set(FlagNames.NoHeader, "true")

	c.Overwrite(set)
	if !c.NoHeader() {
		t.Fatalf("ReadConfig() got = %v, want true", c.NoHeader())
	}
}
