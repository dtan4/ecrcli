package ecr

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/pkg/errors"
)

// Client represents the wrapper of ECR API client
type Client struct {
	api ecriface.ECRAPI
}

// Repository represents the metadata of repository
type Repository struct {
	CreatedAt time.Time
	Name      string
	ARN       string
	URI       string
}

// NewClient creates new Client object
func NewClient(api ecriface.ECRAPI) *Client {
	return &Client{
		api: api,
	}
}

// GetLogin returns ECR login command
func (c *Client) GetLogin() (string, error) {
	resp, err := c.api.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve authorization token")
	}

	if len(resp.AuthorizationData) == 0 {
		return "", errors.New("no authorization data found")
	}

	authData := resp.AuthorizationData[0]

	data, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode authorization data")
	}

	ss := strings.Split(string(data), ":")
	if len(ss) < 2 {
		return "", errors.Errorf("authorization data must be user:pass. got: %q", string(data))
	}

	return fmt.Sprintf("docker login -u %s -p %s -e none %s", ss[0], ss[1], *authData.ProxyEndpoint), nil
}

// ListRepositories returns the list of stored repositories
func (c *Client) ListRepositories() ([]*Repository, error) {
	resp, err := c.api.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		return []*Repository{}, errors.Wrap(err, "failed to retrieve repositories")
	}

	repositories := []*Repository{}

	for _, repository := range resp.Repositories {
		repositories = append(repositories, &Repository{
			CreatedAt: aws.TimeValue(repository.CreatedAt),
			Name:      aws.StringValue(repository.RepositoryName),
			ARN:       aws.StringValue(repository.RepositoryArn),
			URI:       aws.StringValue(repository.RepositoryUri),
		})
	}

	return repositories, nil
}
