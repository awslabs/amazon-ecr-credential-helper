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
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
)

// Options makes the constructors more configurable
type Options struct {
	Session  *session.Session
	Config   *aws.Config
	CacheDir string
}

type ClientFactory interface {
	NewClient(awsSession *session.Session, awsConfig *aws.Config) Client
	NewClientWithOptions(opts Options) Client
	NewClientFromRegion(region string) Client
	NewClientWithDefaults() Client
}
type DefaultClientFactory struct{}

// NewClientWithDefaults creates the client and defaults region
func (defaultClientFactory DefaultClientFactory) NewClientWithDefaults() Client {
	awsSession := session.New()
	awsConfig := awsSession.Config
	return defaultClientFactory.NewClientWithOptions(Options{
		Session: awsSession,
		Config:  awsConfig,
	})
}

// NewClientFromRegion uses the region to create the client
func (defaultClientFactory DefaultClientFactory) NewClientFromRegion(region string) Client {
	awsSession := session.New()
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
	return &defaultClient{
		ecrClient:       ecr.New(opts.Session, opts.Config),
		credentialCache: cache.BuildCredentialsCache(opts.Session, aws.StringValue(opts.Config.Region), opts.CacheDir),
	}
}
