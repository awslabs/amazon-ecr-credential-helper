// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
)

type MockClientFactory struct {
	NewClientFn                 func(ctx context.Context, awsConfig aws.Config) (api.Client, error)
	NewClientWithOptionsFn      func(ctx context.Context, opts api.Options) (api.Client, error)
	NewClientFromRegionFn       func(ctx context.Context, region string) (api.Client, error)
	NewClientWithFipsEndpointFn func(ctx context.Context, region string) (api.Client, error)
	NewClientWithDefaultsFn     func(ctx context.Context) (api.Client, error)
}

func (m MockClientFactory) NewClient(ctx context.Context, awsConfig aws.Config) (api.Client, error) {
	return m.NewClientFn(ctx, awsConfig)
}

func (m MockClientFactory) NewClientWithOptions(ctx context.Context, opts api.Options) (api.Client, error) {
	return m.NewClientWithOptionsFn(ctx, opts)
}

func (m MockClientFactory) NewClientFromRegion(ctx context.Context, region string) (api.Client, error) {
	return m.NewClientFromRegionFn(ctx, region)
}

func (m MockClientFactory) NewClientWithFipsEndpoint(ctx context.Context, region string) (api.Client, error) {
	return m.NewClientWithFipsEndpointFn(ctx, region)
}

func (m MockClientFactory) NewClientWithDefaults(ctx context.Context) (api.Client, error) {
	return m.NewClientWithDefaultsFn(ctx)
}

var _ api.ClientFactory = (*MockClientFactory)(nil)

type MockClient struct {
	GetCredentialsFn             func(ctx context.Context, serverURL string) (*api.Auth, error)
	GetCredentialsByRegistryIDFn func(ctx context.Context, registryID string) (*api.Auth, error)
	ListCredentialsFn            func(ctx context.Context) ([]*api.Auth, error)
}

var _ api.Client = (*MockClient)(nil)

func (m *MockClient) GetCredentials(ctx context.Context, serverURL string) (*api.Auth, error) {
	return m.GetCredentialsFn(ctx, serverURL)
}

func (m *MockClient) GetCredentialsByRegistryID(ctx context.Context, registryID string) (*api.Auth, error) {
	return m.GetCredentialsByRegistryIDFn(ctx, registryID)
}

func (m *MockClient) ListCredentials(ctx context.Context) ([]*api.Auth, error) {
	return m.ListCredentialsFn(ctx)
}
