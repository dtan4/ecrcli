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

	expected := "docker login -u username -p password https://012345678910.dkr.ecr.us-east-1.amazonaws.com"

	got, err := client.GetLogin()
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if got != expected {
		t.Errorf("output does not match. expected: %q, got: %q", expected, got)
	}
}

func TestListImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := "repository"
	pushedAt := time.Unix(1500532805, 0) // 2017-07-20 15:40:05 +0900

	api := mock.NewMockECRAPI(ctrl)
	api.EXPECT().DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(repository),
	}).Return(&ecr.DescribeImagesOutput{
		ImageDetails: []*ecr.ImageDetail{
			&ecr.ImageDetail{
				RegistryId:     aws.String("012345678910"),
				RepositoryName: aws.String("repository"),
				ImageDigest:    aws.String("sha256:6e6810e09a120ebcc3005741c228fecc7f77c513f6565c736370420fbc570bd8"),
				ImageTags: []*string{
					aws.String("latest"),
				},
				ImageSizeInBytes: aws.Int64(186629610),
				ImagePushedAt:    aws.Time(pushedAt),
			},
			&ecr.ImageDetail{
				RegistryId:     aws.String("012345678910"),
				RepositoryName: aws.String("repository"),
				ImageDigest:    aws.String("sha256:b06dd7943a48e1b3ac5a527f0f835eafd3acccdbf508ae4179c1de77617f2310"),
				ImageTags: []*string{
					aws.String("foo"),
					aws.String("bar"),
				},
				ImageSizeInBytes: aws.Int64(178952648),
				ImagePushedAt:    aws.Time(pushedAt),
			},
			&ecr.ImageDetail{
				RegistryId:       aws.String("012345678910"),
				RepositoryName:   aws.String("repository"),
				ImageDigest:      aws.String("sha256:96cfebabbfb81b9e6bf8d03e6d2e0de0a236d429e885a00c68a2a8e17da7cf93"),
				ImageTags:        []*string{},
				ImageSizeInBytes: aws.Int64(186632884),
				ImagePushedAt:    aws.Time(pushedAt),
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	got, err := client.ListImages(repository)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	expected := []*Image{
		&Image{
			Repository: repository,
			Digest:     "sha256:6e6810e09a120ebcc3005741c228fecc7f77c513f6565c736370420fbc570bd8",
			Tags: []string{
				"latest",
			},
			SizeInBytes: 186629610,
			PushedAt:    pushedAt,
		},
		&Image{
			Repository: repository,
			Digest:     "sha256:b06dd7943a48e1b3ac5a527f0f835eafd3acccdbf508ae4179c1de77617f2310",
			Tags: []string{
				"foo",
				"bar",
			},
			SizeInBytes: 178952648,
			PushedAt:    pushedAt,
		},
		&Image{
			Repository:  repository,
			Digest:      "sha256:96cfebabbfb81b9e6bf8d03e6d2e0de0a236d429e885a00c68a2a8e17da7cf93",
			Tags:        []string{},
			SizeInBytes: 186632884,
			PushedAt:    pushedAt,
		},
	}

	for i := range got {
		if !imageEquals(got[i], expected[i]) {
			t.Errorf("expected[%d]:\n%#v, got[%d]:\n%#v", i, expected[i], i, got[i])
		}
	}
}

func imageEquals(a, b *Image) bool {
	return a.Repository == b.Repository &&
		a.Digest == b.Digest &&
		reflect.DeepEqual(a.Tags, b.Tags) &&
		a.SizeInBytes == b.SizeInBytes &&
		a.PushedAt.Equal(b.PushedAt)
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
