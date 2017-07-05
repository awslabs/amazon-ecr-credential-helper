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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ecr "github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/mock"
)

var (
	_ ECRClient = &MockECRClient{}
)

type ECRClient interface {
	GetAuthorizationToken(*ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

type MockECRClient struct {
	mock.Mock
	GetAuthorizationTokenFn func(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

func (m *MockECRClient) GetAuthorizationToken(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	if m.GetAuthorizationTokenFn != nil {
		return m.GetAuthorizationTokenFn(i)
	}

	args := m.Called(i)
	return args.Get(0).(*ecr.GetAuthorizationTokenOutput), args.Error(1)
}

var (
	_ Client = &MockClient{}
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) GetCredentials(serverURL string) (*Auth, error) {
	args := c.Called(serverURL)
	auth := args.Get(0)
	if auth == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Auth), args.Error(1)
}

func (c *MockClient) GetCredentialsByRegistryID(registryID string) (*Auth, error) {
	args := c.Called(registryID)
	auth := args.Get(0)
	if auth == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Auth), args.Error(1)
}

func (c *MockClient) ListCredentials() ([]*Auth, error) {
	args := c.Called()
	auth := args.Get(0)
	if auth == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Auth), args.Error(1)
}

var (
	_ ClientFactory = &MockClientFactory{}
)

type MockClientFactory struct {
	mock.Mock
}

func (cf *MockClientFactory) NewClient(awsSession *session.Session, awsConfig *aws.Config) Client {
	args := cf.Called(awsSession, awsConfig)
	return args.Get(0).(Client)
}

func (cf *MockClientFactory) NewClientWithOptions(opts Options) Client {
	args := cf.Called(opts)
	return args.Get(0).(Client)
}

func (cf *MockClientFactory) NewClientFromRegion(region string) Client {
	args := cf.Called(region)
	return args.Get(0).(Client)
}

func (cf *MockClientFactory) NewClientWithDefaults() Client {
	args := cf.Called()
	return args.Get(0).(Client)
}
