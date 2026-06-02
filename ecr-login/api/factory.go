// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/transport/http"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/version"
)

// Options makes the constructors more configurable
type Options struct {
	Config   aws.Config
	CacheDir string
}

// ClientFactory is a factory for creating clients to interact with ECR
type ClientFactory interface {
	NewClient(ctx context.Context, awsConfig aws.Config) (Client, error)
	NewClientWithOptions(ctx context.Context, opts Options) (Client, error)
	NewClientFromRegion(ctx context.Context, region string) (Client, error)
	NewClientWithFipsEndpoint(ctx context.Context, region string) (Client, error)
	NewClientWithDefaults(ctx context.Context) (Client, error)
}

// DefaultClientFactory is a default implementation of the ClientFactory
type DefaultClientFactory struct{}

var userAgentLoadOption = config.WithAPIOptions([]func(*middleware.Stack) error{
	http.AddHeaderValue("User-Agent", "amazon-ecr-credential-helper/"+version.Version),
})

// NewClientWithDefaults creates the client and defaults region
func (defaultClientFactory DefaultClientFactory) NewClientWithDefaults(ctx context.Context) (Client, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, userAgentLoadOption)
	if err != nil {
		return nil, fmt.Errorf("loading default AWS config: %w", err)
	}

	return defaultClientFactory.NewClientWithOptions(ctx, Options{Config: awsConfig})
}

// NewClientWithFipsEndpoint overrides the default ECR service endpoint in a given region to use the FIPS endpoint
func (defaultClientFactory DefaultClientFactory) NewClientWithFipsEndpoint(ctx context.Context, region string) (Client, error) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		userAgentLoadOption,
		config.WithRegion(region),
		config.WithEndpointDiscovery(aws.EndpointDiscoveryEnabled),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for FIPS endpoint in %s: %w", region, err)
	}

	return defaultClientFactory.NewClientWithOptions(ctx, Options{Config: awsConfig})
}

// NewClientFromRegion uses the region to create the client
func (defaultClientFactory DefaultClientFactory) NewClientFromRegion(ctx context.Context, region string) (Client, error) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		userAgentLoadOption,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for %s: %w", region, err)
	}

	return defaultClientFactory.NewClientWithOptions(ctx, Options{
		Config: awsConfig,
	})
}

// NewClient Create new client with AWS Config
func (defaultClientFactory DefaultClientFactory) NewClient(ctx context.Context, awsConfig aws.Config) (Client, error) {
	return defaultClientFactory.NewClientWithOptions(ctx, Options{Config: awsConfig})
}

// NewClientWithOptions Create new client with Options
func (defaultClientFactory DefaultClientFactory) NewClientWithOptions(ctx context.Context, opts Options) (Client, error) {
	// The ECR Public API is only available in us-east-1 today
	publicConfig := opts.Config.Copy()
	publicConfig.Region = "us-east-1"
	return &defaultClient{
		ecrClient:       NewECRClientWrapper(ecr.NewFromConfig(opts.Config)),
		ecrPublicClient: NewECRPublicClientWrapper(ecrpublic.NewFromConfig(publicConfig)),
		credentialCache: cache.BuildCredentialsCache(ctx, opts.Config, opts.CacheDir),
	}, nil
}
