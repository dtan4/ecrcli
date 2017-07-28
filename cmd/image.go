package cmd

import (
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image <subcommand>",
	Short: "Image related commands",
}

func init() {
	RootCmd.AddCommand(imageCmd)
}
