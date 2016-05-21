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
	"testing"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/mocks"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	region           = "my-region-1"
	registryID       = "123456789012"
	image            = registryID + ".dkr.ecr." + region + ".amazonaws.com/my-image"
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

	factory.EXPECT().NewClient(region).Return(client)
	client.EXPECT().GetCredentials(registryID, image).Return(expectedUsername, expectedPassword, nil)

	username, password, err := helper.Get(image)
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

	factory.EXPECT().NewClient(region).Return(client)
	client.EXPECT().GetCredentials(registryID, image).Return("", "", errors.New("test error"))

	username, password, err := helper.Get(image)
	assert.Equal(t, credentials.ErrCredentialsNotFound, err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}

func TestGetNoMatch(t *testing.T) {
	helper := &ECRHelper{}

	username, password, err := helper.Get("not-ecr-server-url")
	assert.Equal(t, credentials.ErrCredentialsNotFound, err)
	assert.Empty(t, username)
	assert.Empty(t, password)
}
