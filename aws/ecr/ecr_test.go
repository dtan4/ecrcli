package ecr

import (
	"reflect"
	"testing"
	"time"

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

func TestListRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Unix(1500532805, 0) // 2017-07-20 15:40:05 +0900

	api := mock.NewMockECRAPI(ctrl)
	api.EXPECT().DescribeRepositories(&ecr.DescribeRepositoriesInput{}).Return(&ecr.DescribeRepositoriesOutput{
		Repositories: []*ecr.Repository{
			&ecr.Repository{
				RepositoryArn:  aws.String("arn:aws:ecr:us-east-1:012345678910:repository/foo"),
				RegistryId:     aws.String("012345678910"),
				RepositoryName: aws.String("foo"),
				RepositoryUri:  aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/foo"),
				CreatedAt:      aws.Time(createdAt),
			},
			&ecr.Repository{
				RepositoryArn:  aws.String("arn:aws:ecr:us-east-1:012345678910:repository/bar"),
				RegistryId:     aws.String("012345678910"),
				RepositoryName: aws.String("bar"),
				RepositoryUri:  aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/bar"),
				CreatedAt:      aws.Time(createdAt),
			},
			&ecr.Repository{
				RepositoryArn:  aws.String("arn:aws:ecr:us-east-1:012345678910:repository/baz"),
				RegistryId:     aws.String("012345678910"),
				RepositoryName: aws.String("baz"),
				RepositoryUri:  aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/baz"),
				CreatedAt:      aws.Time(createdAt),
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	got, err := client.ListRepositories()
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	expected := []*Repository{
		&Repository{
			CreatedAt: createdAt,
			Name:      "foo",
			ARN:       "arn:aws:ecr:us-east-1:012345678910:repository/foo",
			URI:       "012345678910.dkr.ecr.us-east-1.amazonaws.com/foo",
		},
		&Repository{
			CreatedAt: createdAt,
			Name:      "bar",
			ARN:       "arn:aws:ecr:us-east-1:012345678910:repository/bar",
			URI:       "012345678910.dkr.ecr.us-east-1.amazonaws.com/bar",
		},
		&Repository{
			CreatedAt: createdAt,
			Name:      "baz",
			ARN:       "arn:aws:ecr:us-east-1:012345678910:repository/baz",
			URI:       "012345678910.dkr.ecr.us-east-1.amazonaws.com/baz",
		},
	}

	for i := range got {
		if !repositoryEquals(got[i], expected[i]) {
			t.Errorf("expected[%d]:\n%#v, got[%d]:\n%#v", i, expected[i], i, got[i])
		}
	}
}

func repositoryEquals(a, b *Repository) bool {
	return a.CreatedAt.Equal(b.CreatedAt) && a.ARN == b.ARN && a.Name == b.Name && a.URI == b.URI
}
