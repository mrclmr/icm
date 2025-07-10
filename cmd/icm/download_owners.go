package cmd

import (
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/mrclmr/icm/internal/data"
	"github.com/mrclmr/icm/internal/http"

	"github.com/spf13/cobra"
)

type filePathValue struct {
	homeDir   string
	ownerPath string
}

func (v *filePathValue) String() string {
	return filepath.Join("$HOME", v.ownerPath)
}

func (v *filePathValue) Path() string {
	return filepath.Join(v.homeDir, v.ownerPath)
}

func (v *filePathValue) Set(value string) error {
	dir, _ := path.Split(value)
	if dir != "" {
		_, err := os.Stat(dir)
		if err != nil {
			return err
		}
	}
	v.homeDir = ""
	v.ownerPath = value
	return nil
}

func (*filePathValue) Type() string {
	return "string"
}

func newDownloadOwnersCmd(
	writeOwnersCSVFunc data.WriteOwnersCSVFunc,
	timestampUpdater data.TimestampUpdater,
	ownersDownloader http.OwnersDownloader,
	homeDir string,
	ownerCSVPath string,
) (*cobra.Command, error) {
	filePath := filePathValue{
		homeDir:   homeDir,
		ownerPath: ownerCSVPath,
	}

	downloadOwnersCmd := &cobra.Command{
		Aliases: []string{"update"},
		Use:     "download-owners",
		Short:   "Download information of owners and write CSV to file",
		Long: `Download information of owners and write CSV to file.
Following information is available:

  Owner code
  Company
  City
  Country`,
		Example: `# Overwrite owner.csv file with newest owners
icm download-owners
# Create custom-owner.csv to have additional custom mapping of owner codes
# Use semicolon as a separator. For using double quotes please see existing
# owner.csv file.
echo 'AAA;my company;my city;my country' >> $HOME/.icm/data/custom-owner.csv`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return overwriteOwnersFile(cmd.Context(), writeOwnersCSVFunc, timestampUpdater, ownersDownloader, filePath.Path())
		},
	}
	downloadOwnersCmd.Flags().VarP(&filePath, "output", "o", "output file")

	err := downloadOwnersCmd.MarkFlagFilename("output")
	if err != nil {
		return nil, err
	}

	return downloadOwnersCmd, nil
}

func overwriteOwnersFile(ctx context.Context, writeOwnersCSV data.WriteOwnersCSVFunc, timestampUpdater data.TimestampUpdater, ownersDownloader http.OwnersDownloader, filePath string) error {
	if err := timestampUpdater.Update(); err != nil {
		return err
	}

	owners, err := ownersDownloader.Download(ctx)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	return writeOwnersCSV(owners, file)
}
