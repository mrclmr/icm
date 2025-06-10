package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocCmd(rootCmd *cobra.Command) *cobra.Command {
	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "Documentation commands for man pages and markdown generation",
		Long:  "Documentation commands for man pages and markdown generation.",
	}

	// https://unix.stackexchange.com/questions/3586/what-do-the-numbers-in-a-man-page-mean
	manCmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generate man pages",
		SilenceUsage:          true,
		Hidden:                true,
		DisableFlagsInUseLine: true,
		Example:               "icm doc man . && cat icm.1",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, args []string) error {
			path := args[0]
			err := doc.GenManTree(rootCmd, nil, path)
			if err != nil {
				return err
			}
			return nil
		},
	}

	mdCmd := &cobra.Command{
		Use:                   "markdown",
		Short:                 "Generate markdown",
		SilenceUsage:          true,
		Hidden:                true,
		DisableFlagsInUseLine: true,
		Example:               "icm doc markdown docs/",
		Args:                  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, args[0])
		},
	}

	docCmd.AddCommand(manCmd)
	docCmd.AddCommand(mdCmd)

	return docCmd
}
