package cmd

import (
	"fmt"
	"os"

	"github.com/dtan4/ecrcli/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootOpts = struct {
	debug  bool
	region string
}{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ecrcli",
	Short: "A brief description of your application",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.Initialize(rootOpts.region); err != nil {
			return errors.Wrap(err, "failed to initialize AWS API clients")
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if rootOpts.debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&rootOpts.debug, "debug", false, "Debug mode")
	RootCmd.PersistentFlags().StringVar(&rootOpts.region, "region", "", "AWS region")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
