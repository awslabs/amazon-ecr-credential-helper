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
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	log "github.com/cihub/seelog"
)

const proxyEndpointScheme = "https://"

type Client interface {
	GetCredentials(registry, image string) (*Auth, error)
	ListCredentials() ([]*Auth, error)
}
type defaultClient struct {
	ecrClient       ecriface.ECRAPI
	credentialCache cache.CredentialsCache
}

type Auth struct {
	ProxyEndpoint string
	Username      string
	Password      string
}

// GetCredentials returns username, password, and proxyEndpoint
func (self *defaultClient) GetCredentials(registry, image string) (*Auth, error) {
	log.Debugf("GetCredentials for %s", registry)

	cachedEntry := self.credentialCache.Get(registry)

	if cachedEntry != nil {
		if cachedEntry.IsValid(time.Now()) {
			log.Debugf("Using cached token for %s", registry)
			return extractToken(cachedEntry.AuthorizationToken, cachedEntry.ProxyEndpoint)
		}
		log.Debugf("Cached token is no longer valid. RequestAt: %s, ExpiresAt: %s", cachedEntry.RequestedAt, cachedEntry.ExpiresAt)
	}

	log.Debugf("Calling ECR.GetAuthorizationToken for %s", registry)

	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(registry)},
	}

	output, err := self.ecrClient.GetAuthorizationToken(input)

	if err != nil || output == nil {
		if err == nil {
			err = fmt.Errorf("Missing AuthorizationData in ECR response for %s", registry)
		}

		// if we have a cached token, fall back to avoid failing the request. This may result an expired token
		// being returned, but if there is a 500 or timeout from the service side, we'd like to attempt to re-use an
		// old token. We invalidate tokens prior to their expiration date to help mitigate this scenario.
		if cachedEntry != nil {
			log.Infof("Got error fetching authorization token. Falling back to cached token. Error was: %s", err)
			return extractToken(cachedEntry.AuthorizationToken, cachedEntry.ProxyEndpoint)
		}

		return nil, err
	}

	for _, authData := range output.AuthorizationData {
		if authData.ProxyEndpoint != nil &&
			containsProxyEndpoint(image, aws.StringValue(authData.ProxyEndpoint)) &&
			authData.AuthorizationToken != nil {
			authEntry := cache.AuthEntry{
				AuthorizationToken: aws.StringValue(authData.AuthorizationToken),
				RequestedAt:        time.Now(),
				ExpiresAt:          aws.TimeValue(authData.ExpiresAt),
				ProxyEndpoint:      aws.StringValue(authData.ProxyEndpoint),
			}

			self.credentialCache.Set(registry, &authEntry)
			return extractToken(aws.StringValue(authData.AuthorizationToken), aws.StringValue(authData.ProxyEndpoint))
		}
	}
	return nil, fmt.Errorf("No AuthorizationToken found for %s", registry)
}

func (self *defaultClient) ListCredentials() ([]*Auth, error) {
	log.Debug("Listing credentials from cache")

	auths := []*Auth{}
	for _, authEntry := range self.credentialCache.List() {
		auth, err := extractToken(authEntry.AuthorizationToken, authEntry.ProxyEndpoint)
		if err != nil {
			log.Debugf("Could not extract token: %v:", err)
		} else {
			auths = append(auths, auth)
		}
	}

	return auths, nil
}

func containsProxyEndpoint(image string, proxyEndpoint string) bool {
	if image == "" {
		return true
	}
	return strings.HasPrefix(proxyEndpointScheme+image, proxyEndpoint)
}

func extractToken(token string, proxyEndpoint string) (*Auth, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v:", err)
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid token: expected two parts, got %n", len(parts))
	}

	return &Auth{
		Username:      parts[0],
		Password:      parts[1],
		ProxyEndpoint: proxyEndpoint,
	}, nil
}
