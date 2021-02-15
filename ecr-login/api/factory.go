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
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/version"
)

// Options makes the constructors more configurable
type Options struct {
	Session  *session.Session
	Config   *aws.Config
	CacheDir string
}

// ClientFactory is a factory for creating clients to interact with ECR
type ClientFactory interface {
	NewClient(awsSession *session.Session, awsConfig *aws.Config) Client
	NewClientWithOptions(opts Options) Client
	NewClientFromRegion(region string) Client
	NewClientWithFipsEndpoint(region string) (Client, error)
	NewClientWithDefaults() Client
}

// DefaultClientFactory is a default implementation of the ClientFactory
type DefaultClientFactory struct{}

var userAgentHandler = request.NamedHandler{
	Name: "ecr-login.UserAgentHandler",
	Fn:   request.MakeAddToUserAgentHandler("amazon-ecr-credential-helper", version.Version),
}

// NewClientWithDefaults creates the client and defaults region
func (defaultClientFactory DefaultClientFactory) NewClientWithDefaults() Client {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: loadSharedConfigState(),
	}))

	awsSession.Handlers.Build.PushBackNamed(userAgentHandler)
	awsConfig := awsSession.Config
	return defaultClientFactory.NewClientWithOptions(Options{
		Session: awsSession,
		Config:  awsConfig,
	})
}

// NewClientWithFipsEndpoint overrides the default ECR service endpoint in a given region to use the FIPS endpoint
func (defaultClientFactory DefaultClientFactory) NewClientWithFipsEndpoint(region string) (Client, error) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: loadSharedConfigState(),
	}))

	awsSession.Handlers.Build.PushBackNamed(userAgentHandler)

	endpoint, err := getServiceEndpoint("ecr-fips", region)
	if err != nil {
		return nil, err
	}

	awsConfig := awsSession.Config.WithEndpoint(endpoint).WithRegion(region)
	return defaultClientFactory.NewClientWithOptions(Options{
		Session: awsSession,
		Config:  awsConfig,
	}), nil
}

// NewClientFromRegion uses the region to create the client
func (defaultClientFactory DefaultClientFactory) NewClientFromRegion(region string) Client {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: loadSharedConfigState(),
	}))
	awsSession.Handlers.Build.PushBackNamed(userAgentHandler)
	awsConfig := &aws.Config{Region: aws.String(region)}
	return defaultClientFactory.NewClientWithOptions(Options{
		Session: awsSession,
		Config:  awsConfig,
	})
}

// NewClient Create new client with AWS Config
func (defaultClientFactory DefaultClientFactory) NewClient(awsSession *session.Session, awsConfig *aws.Config) Client {
	return defaultClientFactory.NewClientWithOptions(Options{
		Session: awsSession,
		Config:  awsConfig,
	})
}

// NewClientWithOptions Create new client with Options
func (defaultClientFactory DefaultClientFactory) NewClientWithOptions(opts Options) Client {
	// The ECR Public API is only available in us-east-1 today
	publicConfig := opts.Config.Copy().WithRegion("us-east-1")
	return &defaultClient{
		ecrClient:       ecr.New(opts.Session, opts.Config),
		ecrPublicClient: ecrpublic.New(opts.Session, publicConfig),
		credentialCache: cache.BuildCredentialsCache(opts.Session, aws.StringValue(opts.Config.Region), opts.CacheDir),
	}
}

func getServiceEndpoint(service, region string) (string, error) {
	resolver := endpoints.DefaultResolver()
	endpoint, err := resolver.EndpointFor(service, region, func(opts *endpoints.Options) {
		opts.ResolveUnknownService = true
	})
	return endpoint.URL, err
}

func loadSharedConfigState() session.SharedConfigState {
	loadConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")
	if loadConfig == "" {
		return session.SharedConfigEnable
	}
	if enable, err := strconv.ParseBool(loadConfig); err == nil && !enable {
		return session.SharedConfigDisable
	}
	return session.SharedConfigEnable
}
