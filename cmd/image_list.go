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
	imageListHeader = []string{
		"DIGEST",
		"PUSHEDAT",
		"TAGS",
	}
)

// imageListCmd represents the imageList command
var imageListCmd = &cobra.Command{
	Use:   "list REPO",
	Short: "List images",
	RunE:  doImageList,
}

func doImageList(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("repository name must be given")
	}
	repo := args[0]

	images, err := aws.ECR.ListImages(repo)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch image list of %s", repo)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(imageListHeader, "\t"))

	for _, image := range images {
		fmt.Fprintln(w, strings.Join([]string{
			image.Digest,
			image.PushedAt.Local().String(),
			strings.Join(image.Tags, ","),
		}, "\t"))
	}

	w.Flush()

	return nil
}

func init() {
	imageCmd.AddCommand(imageListCmd)
}
