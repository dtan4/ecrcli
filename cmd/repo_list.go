package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/ecrcli/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	repoListHeader = []string{
		"NAME",
		"URI",
		"CREATEDAT",
	}
)

// repoListCmd represents the repoList command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	RunE:  doRepoList,
}

func doRepoList(cmd *cobra.Command, args []string) error {
	repos, err := aws.ECR.ListRepositories()
	if err != nil {
		return errors.Wrap(err, "failed to fetch repository list")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(repoListHeader, "\t"))

	for _, repo := range repos {
		fmt.Fprintln(w, strings.Join([]string{
			repo.Name,
			repo.URI,
			repo.CreatedAt.Local().String(),
		}, "\t"))
	}

	w.Flush()

	return nil
}

func init() {
	repoCmd.AddCommand(repoListCmd)
}
