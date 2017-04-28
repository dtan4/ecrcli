package ecr

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/dtan4/ecrcli/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestGetLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockECRAPI(ctrl)
	api.EXPECT().GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				AuthorizationToken: aws.String("dXNlcm5hbWU6cGFzc3dvcmQ="),
				ProxyEndpoint:      aws.String("https://012345678910.dkr.ecr.us-east-1.amazonaws.com"),
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := "docker login -u username -p password -e none https://012345678910.dkr.ecr.us-east-1.amazonaws.com"

	got, err := client.GetLogin()
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if got != expected {
		t.Errorf("output does not match. expected: %q, got: %q", expected, got)
	}
}
