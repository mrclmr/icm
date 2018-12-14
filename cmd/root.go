// Copyright © 2017 Marcel Meyer meyermarcel@posteo.de
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
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/meyermarcel/icm/internal/data"

	"github.com/meyermarcel/icm/configs"

	"path/filepath"

	"os/user"

	"github.com/meyermarcel/icm/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type decoders struct {
	ownerDecodeUpdater data.OwnerDecodeUpdater
	equipCatDecoder    data.EquipCatDecoder
	sizeTypeDecoders
}

type sizeTypeDecoders struct {
	lengthDecoder      data.LengthDecoder
	heightWidthDecoder data.HeightWidthDecoder
	typeDecoder        data.TypeDecoder
}

const (
	appName  = "icm"
	appDir   = "." + appName
	ownerURL = "https://www.bic-code.org/bic-letter-search/?resultsperpage=17576&searchterm="
)

var sepHelp = `Configuration for separators is generated first time you
execute a command that requires the configuration.

Flags for output formatting can overridden with a config file.
Edit default configuration:

  ` + filepath.Join("$HOME", appDir, configs.NameWithYmlExt)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	stderr := os.Stderr

	current, err := user.Current()
	checkErr(stderr, err)

	appDirPath := initDir(filepath.Join(current.HomeDir, appDir))

	pathToCfg := filepath.Join(appDirPath, configs.NameWithYmlExt)
	if _, err := os.Stat(pathToCfg); os.IsNotExist(err) {
		err = ioutil.WriteFile(pathToCfg, configs.Cfg(), 0644)
		checkErr(stderr, err)
	}
	viperCfg := viper.New()
	viperCfg.AddConfigPath(appDirPath)
	viperCfg.SetConfigName(configs.Name)
	err = viperCfg.ReadInConfig()
	checkErr(stderr, err)

	appDirDataPath := initDir(filepath.Join(appDirPath, "data"))

	ownerDecodeUpdater, err := file.NewOwnerDecoderUpdater(appDirDataPath)
	checkErr(stderr, err)

	equipCatDecoder, err := file.NewEquipCatDecoder(appDirDataPath)
	checkErr(stderr, err)

	lengthDecoder, heightWidthDecoder, err := file.NewSizeDecoder(appDirDataPath)
	checkErr(stderr, err)

	typeDecoder, err := file.NewTypeDecoder(appDirDataPath)
	checkErr(stderr, err)

	timestampUpdater, err := file.NewTimestampUpdater(appDirDataPath)
	checkErr(stderr, err)

	rootCmd := newRootCmd(
		version,
		os.Stdout,
		stderr,
		viperCfg,
		decoders{
			ownerDecodeUpdater,
			equipCatDecoder,
			sizeTypeDecoders{
				lengthDecoder,
				heightWidthDecoder,
				typeDecoder},
		},
		timestampUpdater,
		ownerURL)

	err = rootCmd.Execute()
	checkCmdErr(err)
	checkErr(stderr, err)
}

func newRootCmd(
	version string,
	writer, writerErr io.Writer,
	viper *viper.Viper,
	decoders decoders,
	timestampUpdater data.TimestampUpdater,
	ownerURL string) *cobra.Command {
	rootCmd := &cobra.Command{
		Version:       version,
		Use:           appName,
		Short:         "Validate or generate intermodal container markings",
		Long:          "Validate or generate intermodal container markings.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newCompletionCmd(writer, rootCmd))
	rootCmd.AddCommand(newGenerateCmd(writer, writerErr, viper, decoders.ownerDecodeUpdater))
	rootCmd.AddCommand(newValidateCmd(writer, writerErr, viper, decoders))
	rootCmd.AddCommand(newUpdateOwnerCmd(decoders.ownerDecodeUpdater, timestampUpdater, ownerURL))

	return rootCmd
}

func initDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModeDir|0700)
	}
	return path
}

func checkCmdErr(err error) {
	switch err.(type) {
	case *errValidate:
		os.Exit(1)
	}
}

func writeErr(writer io.Writer, err error) {
	if err != nil {
		_, _ = fmt.Fprintf(writer, "%s: %s\n", appName, err)
	}
}

func checkErr(writer io.Writer, err error) {
	if err != nil {
		writeErr(writer, err)
		os.Exit(1)
	}
}
