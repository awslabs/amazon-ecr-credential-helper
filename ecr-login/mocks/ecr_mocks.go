// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package mock_api

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
)

type MockClientFactory struct {
	NewClientFn                 func(awsConfig *aws.Config) api.Client
	NewClientWithOptionsFn      func(opts api.Options) api.Client
	NewClientFromRegionFn       func(region string) api.Client
	NewClientWithFipsEndpointFn func(region string) (api.Client, error)
	NewClientWithDefaultsFn     func() api.Client
}

func (m MockClientFactory) NewClient(awsConfig *aws.Config) api.Client {
	return m.NewClientFn(awsConfig)
}

func (m MockClientFactory) NewClientWithOptions(opts api.Options) api.Client {
	return m.NewClientWithOptionsFn(opts)
}

func (m MockClientFactory) NewClientFromRegion(region string) api.Client {
	return m.NewClientFromRegionFn(region)
}

func (m MockClientFactory) NewClientWithFipsEndpoint(region string) (api.Client, error) {
	return m.NewClientWithFipsEndpointFn(region)
}

func (m MockClientFactory) NewClientWithDefaults() api.Client {
	return m.NewClientWithDefaultsFn()
}

var _ api.ClientFactory = (*MockClientFactory)(nil)

type MockClient struct {
	GetCredentialsFn             func(serverURL string) (*api.Auth, error)
	GetCredentialsByRegistryIDFn func(registryID string) (*api.Auth, error)
	ListCredentialsFn            func() ([]*api.Auth, error)
}

var _ api.Client = (*MockClient)(nil)

func (m *MockClient) GetCredentials(serverURL string) (*api.Auth, error) {
	return m.GetCredentialsFn(serverURL)
}

func (m *MockClient) GetCredentialsByRegistryID(registryID string) (*api.Auth, error) {
	return m.GetCredentialsByRegistryIDFn(registryID)
}

func (m *MockClient) ListCredentials() ([]*api.Auth, error) {
	return m.ListCredentialsFn()
}
