package ecr

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/pkg/errors"
)

// Client represents the wrapper of ECR API client
type Client struct {
	api ecriface.ECRAPI
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
