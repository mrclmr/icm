package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/meyermarcel/icm/data/file"

	"github.com/meyermarcel/icm/data"

	"github.com/meyermarcel/icm/configs"

	"github.com/spf13/cobra"
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
	ownerURL = "https://www.bic-code.org/search/bic-codes/country/all/results/17576"
)

var sepHelp = `Configuration for separators is generated first time you
execute a command that requires the configuration.

Flags for output formatting can be overridden with a config file.
Edit default configuration for customization:

  ` + filepath.Join("$HOME", appDir, configs.ConfigNameWithYmlExt)

func Execute(version string) {
	stderr := os.Stderr

	homeDir, err := os.UserHomeDir()
	checkErr(stderr, err)

	appDirPath := initDir(filepath.Join(homeDir, appDir))

	pathToCfg := filepath.Join(appDirPath, configs.ConfigNameWithYmlExt)
	if _, err := os.Stat(pathToCfg); os.IsNotExist(err) {
		errWrite := os.WriteFile(pathToCfg, configs.DefaultConfig(), 0o644)
		checkErr(stderr, errWrite)
	}

	cfgFile, err := os.ReadFile(pathToCfg)
	checkErr(stderr, err)
	config, err := configs.ReadConfig(cfgFile)
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

	bufWriter := bufio.NewWriter(os.Stdout)
	rootCmd := newRootCmd(
		version,
		bufWriter,
		stderr,
		config,
		decoders{
			ownerDecodeUpdater,
			equipCatDecoder,
			sizeTypeDecoders{
				lengthDecoder,
				heightWidthDecoder,
				typeDecoder,
			},
		},
		timestampUpdater,
		ownerURL)

	errCmd := rootCmd.Execute()

	errBuf := bufWriter.Flush()
	writeErr(stderr, errBuf)

	var errValidate *validateError
	if errors.As(err, &errValidate) {
		os.Exit(1)
	}

	checkErr(stderr, errCmd)
}

func newRootCmd(
	version string,
	writer, writerErr io.Writer,
	config *configs.Config,
	decoders decoders,
	timestampUpdater data.TimestampUpdater,
	ownerURL string,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Version:           version,
		Use:               appName,
		Short:             "Validate or generate intermodal container markings",
		Long:              "Validate or generate intermodal container markings.",
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Visit github.com/meyermarcel/icm for more docs, issues, pull requests and feedback.
`)

	rootCmd.AddCommand(newGenerateCmd(writer, writerErr, config, decoders.ownerDecodeUpdater))
	rootCmd.AddCommand(newValidateCmd(os.Stdin, writer, config, decoders))
	rootCmd.AddCommand(newUpdateOwnerCmd(decoders.ownerDecodeUpdater, timestampUpdater, ownerURL))
	rootCmd.AddCommand(newDocCmd(rootCmd))

	return rootCmd
}

func initDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModeDir|0o700)
	}
	return path
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
