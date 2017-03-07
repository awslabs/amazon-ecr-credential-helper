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
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	log "github.com/cihub/seelog"
)

const proxyEndpointScheme = "https://"
const programName = "docker-credential-ecr-login"

var ecrPattern = regexp.MustCompile(`(^[a-zA-Z0-9][a-zA-Z0-9-_]*)\.dkr\.ecr\.([a-zA-Z0-9][a-zA-Z0-9-_]*)\.amazonaws\.com(\.cn)?`)

type Registry struct {
	ID	string
	Region	string
}

func ExtractRegistry(serverURL string) (*Registry, error) {
	if (strings.HasPrefix(serverURL, proxyEndpointScheme)) {
		serverURL = strings.TrimPrefix(serverURL, proxyEndpointScheme)
	}
	matches := ecrPattern.FindStringSubmatch(serverURL)
	if len(matches) == 0 {
		return nil, fmt.Errorf(programName + " can only be used with Amazon EC2 Container Registry.")
	} else if len(matches) < 3 {
		return nil, fmt.Errorf(serverURL + "is not a valid repository URI for Amazon EC2 Container Registry.")
	}
	registry := &Registry{
		ID:	matches[1],
	        Region:	matches[2],
	}
	return registry, nil
}

type Client interface {
	GetCredentials(serverURL string) (*Auth, error)
	GetCredentialsByRegistryID(registryID string) (*Auth, error)
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
func (self *defaultClient) GetCredentials(serverURL string) (*Auth, error) {
	registry, err := ExtractRegistry(serverURL)
	if err != nil {
		return nil, err
	}
	log.Debugf("Retrieving credentials for %s in %s (%s)", registry.ID, registry.Region, serverURL)
	return self.GetCredentialsByRegistryID(registry.ID)
}

// GetCredentials returns username, password, and proxyEndpoint
func (self *defaultClient) GetCredentialsByRegistryID(registryID string) (*Auth, error) {
	cachedEntry := self.credentialCache.Get(registryID)
	if cachedEntry != nil {
		if cachedEntry.IsValid(time.Now()) {
			log.Debugf("Using cached token for %s", registryID)
			return extractToken(cachedEntry.AuthorizationToken, cachedEntry.ProxyEndpoint)
		}
		log.Debugf("Cached token is no longer valid. RequestAt: %s, ExpiresAt: %s", cachedEntry.RequestedAt, cachedEntry.ExpiresAt)
	}

	auth, err := self.getAuthorizationToken(registryID)

	// if we have a cached token, fall back to avoid failing the request. This may result an expired token
	// being returned, but if there is a 500 or timeout from the service side, we'd like to attempt to re-use an
	// old token. We invalidate tokens prior to their expiration date to help mitigate this scenario.
	if err != nil && cachedEntry != nil {
		log.Infof("Got error fetching authorization token. Falling back to cached token. Error was: %s", err)
		return extractToken(cachedEntry.AuthorizationToken, cachedEntry.ProxyEndpoint)
	}
	return auth, err
}

func (self *defaultClient) ListCredentials() ([]*Auth, error) {
	auths := []*Auth{}
	for _, authEntry := range self.credentialCache.List() {
		auth, err := extractToken(authEntry.AuthorizationToken, authEntry.ProxyEndpoint)
		if err != nil {
			log.Debugf("Could not extract token: %v", err)
		} else {
			auths = append(auths, auth)
		}
	}

	// If cache is empty, get authorization token of default registry
	if len(auths) == 0 {
		log.Debug("No credential cache")
		auth, err := self.getAuthorizationToken("")
		if err != nil {
			log.Debugf("Couldn't get authorization token: %v", err)
		} else {
			auths = append(auths, auth)
		}
		return auths, err
	}

	return auths, nil 
}

func (self *defaultClient) getAuthorizationToken(registryID string) (*Auth, error) {
	var input *ecr.GetAuthorizationTokenInput
        if registryID == "" {
		log.Debugf("Calling ECR.GetAuthorizationToken for default registry")
		input = &ecr.GetAuthorizationTokenInput{}
	} else {
		log.Debugf("Calling ECR.GetAuthorizationToken for %s", registryID)
		input = &ecr.GetAuthorizationTokenInput{
			RegistryIds: []*string{aws.String(registryID)},
		}
	}

	output, err := self.ecrClient.GetAuthorizationToken(input)
	if err != nil || output == nil {
		if err == nil {
			if registryID == "" {
				err = fmt.Errorf("Mising AuthorizationData in ECR response for default registry")
			} else {
				err = fmt.Errorf("Missing AuthorizationData in ECR response for %s", registryID)
			}
		}
		return nil, err
	}

	for _, authData := range output.AuthorizationData {
		if authData.ProxyEndpoint != nil && authData.AuthorizationToken != nil {
			authEntry := cache.AuthEntry{
				AuthorizationToken: aws.StringValue(authData.AuthorizationToken),
				RequestedAt:        time.Now(),
				ExpiresAt:          aws.TimeValue(authData.ExpiresAt),
				ProxyEndpoint:      aws.StringValue(authData.ProxyEndpoint),
			}
			registry, err := ExtractRegistry(authEntry.ProxyEndpoint)
			if err != nil {
				return nil, fmt.Errorf("Invalid ProxyEndpoint returned by ECR: %s", authEntry.ProxyEndpoint)
			}
			auth, err := extractToken(authEntry.AuthorizationToken, authEntry.ProxyEndpoint)
			if err != nil {
				return nil, err
			}
			self.credentialCache.Set(registry.ID, &authEntry)
			return auth, nil
		}
	}
	if registryID == "" {
		return nil, fmt.Errorf("No AuthorizationToken found for default registry")
	} else {
		return nil, fmt.Errorf("No AuthorizationToken found for %s", registryID)
	}
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
