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

package ecr

import (
	"errors"
	"fmt"
	"testing"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/mocks"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	region           = "us-east-1"
	proxyEndpoint    = "123456789012" + ".dkr.ecr." + region + ".amazonaws.com"
	proxyEndpointUrl = "https://" + proxyEndpoint
	expectedUsername = "username"
	expectedPassword = "password"
)

func TestGetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := mock_api.NewMockClientFactory(ctrl)
	client := mock_api.NewMockClient(ctrl)

	helper := &ECRHelper{
		ClientFactory: factory,
	}

	factory.EXPECT().NewClientFromRegion(region).Return(client)
	client.EXPECT().GetCredentials(proxyEndpoint).Return(&ecr.Auth{
		Username:      expectedUsername,
		Password:      expectedPassword,
		ProxyEndpoint: proxyEndpointUrl,
	}, nil)

	username, password, err := helper.Get(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, expectedUsername, username)
	assert.Equal(t, expectedPassword, password)
}

func TestGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := mock_api.NewMockClientFactory(ctrl)
	client := mock_api.NewMockClient(ctrl)

	helper := &ECRHelper{
		ClientFactory: factory,
	}

	factory.EXPECT().NewClientFromRegion(region).Return(client)
	client.EXPECT().GetCredentials(proxyEndpoint).Return(nil, errors.New("test error"))

	username, password, err := helper.Get(proxyEndpoint)
	assert.True(t, credentials.IsErrCredentialsNotFound(err))
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestGetNoMatch(t *testing.T) {
	helper := &ECRHelper{}

	username, password, err := helper.Get("not-ecr-server-url")
	assert.True(t, credentials.IsErrCredentialsNotFound(err))
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := mock_api.NewMockClientFactory(ctrl)
	client := mock_api.NewMockClient(ctrl)

	helper := &ECRHelper{
		ClientFactory: factory,
	}

	factory.EXPECT().NewClientWithDefaults().Return(client)
	client.EXPECT().ListCredentials().Return([]*ecr.Auth{
		&ecr.Auth{
			Username:      expectedUsername,
			Password:      expectedPassword,
			ProxyEndpoint: proxyEndpointUrl,
		},
	}, nil)

	serverList, err := helper.List()
	assert.NoError(t, err)
	assert.Len(t, serverList, 1)
	assert.Equal(t, expectedUsername, serverList[proxyEndpointUrl])
}

func TestListFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	factory := mock_api.NewMockClientFactory(ctrl)
	client := mock_api.NewMockClient(ctrl)

	helper := &ECRHelper{
		ClientFactory: factory,
	}

	factory.EXPECT().NewClientWithDefaults().Return(client)
	client.EXPECT().ListCredentials().Return(nil, fmt.Errorf("nope"))

	serverList, err := helper.List()
	assert.Error(t, err)
	assert.Len(t, serverList, 0)
}
