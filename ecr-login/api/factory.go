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
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	"github.com/mitchellh/go-homedir"

	log "github.com/cihub/seelog"
)

type ClientFactory interface {
	NewClient(awsSession *session.Session, awsConfig *aws.Config) Client
	NewClientWithRegion(region string) Client
}
type DefaultClientFactory struct{}

// NewClientWithRegion uses the region to create the client
func (defaultClientFactory DefaultClientFactory) NewClientWithRegion(region string) Client {
	awsSession := session.New()
	awsConfig := &aws.Config{Region: aws.String(region)}

	return defaultClientFactory.NewClient(awsSession, awsConfig)
}

// NewClient Create new client with AWS Config
func (defaultClientFactory DefaultClientFactory) NewClient(awsSession *session.Session, awsConfig *aws.Config) Client {
	return &defaultClient{
		ecrClient:       ecr.New(awsSession, awsConfig),
		credentialCache: defaultClientFactory.buildCredentialsCache(awsSession, aws.StringValue(awsConfig.Region)),
	}
}

func (defaultClientFactory DefaultClientFactory) buildCredentialsCache(awsSession *session.Session, region string) cache.CredentialsCache {
	if os.Getenv("AWS_ECR_DISABLE_CACHE") != "" {
		log.Debug("Cache disabled due to AWS_ECR_DISABLE_CACHE")
		return cache.NewNullCredentialsCache()
	}

	cacheDir, err := homedir.Expand("~/.ecr")
	if err != nil {
		log.Debugf("Could expand cache path: %s", err)
		log.Debug("Disabling cache")
		return cache.NewNullCredentialsCache()
	}

	cacheFilename := "cache.json"

	credentials, err := awsSession.Config.Credentials.Get()
	if err != nil {
		log.Debugf("Could fetch credentials for cache prefix: %s", err)
		log.Debug("Disabling cache")
		return cache.NewNullCredentialsCache()
	}

	return cache.NewFileCredentialsCache(cacheDir, cacheFilename, defaultClientFactory.credentialsCachePrefix(region, &credentials))
}

// Determine a key prefix for a credentials cache. Because auth tokens are scoped to an account and region, rely on provided
// region, as well as hash of the access key.
func (defaultClientFactory DefaultClientFactory) credentialsCachePrefix(region string, credentials *credentials.Value) string {
	return fmt.Sprintf("%s-%s-", region, checksum(credentials.AccessKeyID))
}

// Base64 encodes an MD5 checksum. Relied on for uniqueness, and not for cryptographic security.
func checksum(text string) string {
	hasher := md5.New()
	data := hasher.Sum([]byte(text))
	return base64.StdEncoding.EncodeToString(data)
}
