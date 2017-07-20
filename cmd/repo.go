package cmd

import (
	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo <subcommand>",
	Short: "Repository related commands",
}

func init() {
	RootCmd.AddCommand(repoCmd)
}
