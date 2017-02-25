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

type ClientFactory interface {
	NewClient(awsSession *session.Session, awsConfig *aws.Config) Client
	NewClientFromRegion(region string) Client
	NewClientWithDefaults() Client
}
type DefaultClientFactory struct{}

// NewClientWithDefaults creates the client and defaults region
func (defaultClientFactory DefaultClientFactory) NewClientWithDefaults() Client {
	awsSession := session.New()
	return defaultClientFactory.NewClient(awsSession, awsSession.Config)
}

// NewClientFromRegion uses the region to create the client
func (defaultClientFactory DefaultClientFactory) NewClientFromRegion(region string) Client {
	awsSession := session.New()
	awsConfig := &aws.Config{Region: aws.String(region)}

	return defaultClientFactory.NewClient(awsSession, awsConfig)
}

// NewClient Create new client with AWS Config
func (defaultClientFactory DefaultClientFactory) NewClient(awsSession *session.Session, awsConfig *aws.Config) Client {
	return &defaultClient{
		ecrClient:       ecr.New(awsSession, awsConfig),
		credentialCache: cache.BuildCredentialsCache(awsSession, aws.StringValue(awsConfig.Region)),
	}
}
