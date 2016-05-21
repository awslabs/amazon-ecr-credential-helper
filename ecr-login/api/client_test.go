// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package api

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	registryID       = "0123456789012"
	proxyEndpoint    = "proxy"
	expectedUsername = "username"
	expectedPassword = "password"
)

func TestGetAuthConfigSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)

	client := &defaultClient{
		ecrClient: ecrClient,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + proxyEndpoint),
				AuthorizationToken: aws.String(base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))),
			},
		},
	}, nil)

	username, password, err := client.GetCredentials(registryID, proxyEndpoint+"/myimage")
	assert.Nil(t, err)
	assert.Equal(t, username, expectedUsername)
	assert.Equal(t, password, expectedPassword)
}

func TestGetAuthConfigNoMatchAuthorizationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)

	client := &defaultClient{
		ecrClient: ecrClient,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + "notproxy"),
				AuthorizationToken: aws.String(base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))),
			},
		},
	}, nil)

	username, password, err := client.GetCredentials(registryID, proxyEndpoint+"/myimage")
	assert.NotNil(t, err)
	t.Log(err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestGetAuthConfigBadBase64(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)

	client := &defaultClient{
		ecrClient: ecrClient,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(proxyEndpoint),
				AuthorizationToken: aws.String(expectedUsername + ":" + expectedPassword),
			},
		},
	}, nil)

	username, password, err := client.GetCredentials(registryID, proxyEndpoint+"/myimage")
	assert.NotNil(t, err)
	t.Log(err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestGetAuthConfigMissingResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)

	client := &defaultClient{
		ecrClient: ecrClient,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		})

	username, password, err := client.GetCredentials(registryID, proxyEndpoint+"/myimage")
	assert.NotNil(t, err)
	t.Log(err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestGetAuthConfigECRError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)

	client := &defaultClient{
		ecrClient: ecrClient,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(nil, errors.New("test error"))

	username, password, err := client.GetCredentials(registryID, proxyEndpoint+"/myimage")
	assert.NotNil(t, err)
	t.Log(err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}
