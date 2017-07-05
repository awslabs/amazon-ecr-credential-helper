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
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ecr "github.com/aws/aws-sdk-go/service/ecr"
)

var (
	_ ECRClient = &MockECRClient{}
)

type ECRClient interface {
	GetAuthorizationToken(*ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

type MockECRClient struct {
	GetAuthorizationTokenFn func(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

func (m *MockECRClient) GetAuthorizationToken(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	if m.GetAuthorizationTokenFn != nil {
		return m.GetAuthorizationTokenFn(i)
	}
	return nil, nil
}

var (
	_ Client = &MockClient{}
)

type MockClient struct {
	GetCredentialsFn             func(serverURL string) (*Auth, error)
	GetCredentialsByRegistryIDFn func(registryID string) (*Auth, error)
	ListCredentialsFn            func() ([]*Auth, error)
}

func (c *MockClient) GetCredentials(serverURL string) (*Auth, error) {
	if c.GetCredentialsFn != nil {
		return c.GetCredentialsFn(serverURL)
	}
	return nil, errors.New("No mocked function")
}

func (c *MockClient) GetCredentialsByRegistryID(registryID string) (*Auth, error) {
	if c.GetCredentialsByRegistryIDFn != nil {
		return c.GetCredentialsByRegistryIDFn(registryID)
	}
	return nil, errors.New("No mocked function")
}

func (c *MockClient) ListCredentials() ([]*Auth, error) {
	if c.ListCredentialsFn != nil {
		return c.ListCredentialsFn()
	}
	return nil, errors.New("No mocked function")
}

var (
	_ ClientFactory = &MockClientFactory{}
)

type MockClientFactory struct {
	NewClientFn             func(awsSession *session.Session, awsConfig *aws.Config) Client
	NewClientWithOptionsFn  func(opts Options) Client
	NewClientFromRegionFn   func(region string) Client
	NewClientWithDefaultsFn func() Client
}

func (cf *MockClientFactory) NewClient(awsSession *session.Session, awsConfig *aws.Config) Client {
	if cf.NewClientFn != nil {
		return cf.NewClientFn(awsSession, awsConfig)
	}
	return nil
}

func (cf *MockClientFactory) NewClientWithOptions(opts Options) Client {
	if cf.NewClientWithOptionsFn != nil {
		return cf.NewClientWithOptionsFn(opts)
	}
	return nil
}

func (cf *MockClientFactory) NewClientFromRegion(region string) Client {
	if cf.NewClientFromRegionFn != nil {
		return cf.NewClientFromRegionFn(region)
	}
	return nil
}

func (cf *MockClientFactory) NewClientWithDefaults() Client {
	if cf.NewClientWithDefaultsFn != nil {
		return cf.NewClientWithDefaultsFn()
	}
	return nil
}
