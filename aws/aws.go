package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ecrapi "github.com/aws/aws-sdk-go/service/ecr"
	"github.com/dtan4/ecrcli/aws/ecr"
	"github.com/pkg/errors"
)

var (
	// ECR represents ECR API client
	ECR *ecr.Client
)

// Initialize creates AWS API clients
func Initialize(region string) error {
	var sess *session.Session

	if region == "" {
		s, err := session.NewSession()
		if err != nil {
			return errors.Wrap(err, "failed to create new AWS session")
		}

		sess = s
	} else {
		s, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create new AWS session")
		}

		sess = s
	}

	ECR = ecr.NewClient(ecrapi.New(sess))

	return nil
}
