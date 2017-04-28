package cmd

import (
	"fmt"

	"github.com/dtan4/ecrcli/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// getLoginCmd represents the get-login command
var getLoginCmd = &cobra.Command{
	Use:   "get-login",
	Short: "Print ECR login command",
	RunE:  doGetLogin,
}

func doGetLogin(cmd *cobra.Command, args []string) error {
	loginCmd, err := aws.ECR.GetLogin()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve docker login command")
	}

	fmt.Println(loginCmd)

	return nil
}

func init() {
	RootCmd.AddCommand(getLoginCmd)
}
