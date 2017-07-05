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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	"github.com/stretchr/testify/assert"
)

const (
	registryID       = "123456789012"
	proxyEndpoint    = "123456789012.dkr.ecr.us-east-1.amazonaws.com"
	expectedUsername = "username"
	expectedPassword = "password"
)

func TestGetAuthConfigSuccess(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(testProxyEndpoint),
					ExpiresAt:          aws.Time(expiresAt),
					AuthorizationToken: aws.String(authorizationToken),
				},
			},
		}, nil
	}

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
	}

	credentialCache.On("Get", registryID).Return(nil)
	credentialCache.SetFn = func(registry string, actual *cache.AuthEntry) {
		compareAuthEntry(t, actual, authEntry)
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigNoMatchAuthorizationToken(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(proxyEndpointScheme + "notproxy"),
					AuthorizationToken: aws.String(base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))),
				},
			},
		}, nil
	}

	credentialCache.On("Get", registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigGetCacheSuccess(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		ExpiresAt:          expiresAt,
		RequestedAt:        time.Now(),
		AuthorizationToken: authorizationToken,
	}

	credentialCache.On("Get", registryID).Return(authEntry)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigSuccessInvalidCacheHit(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(testProxyEndpoint),
					ExpiresAt:          aws.Time(expiresAt),
					AuthorizationToken: aws.String(authorizationToken),
				},
			},
		}, nil
	}

	expiredAuthEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now().Add(-12 * time.Hour),
		ExpiresAt:          time.Now().Add(-6 * time.Hour),
		AuthorizationToken: authorizationToken,
	}

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
	}

	credentialCache.On("Get", registryID).Return(expiredAuthEntry)
	credentialCache.SetFn = func(registry string, actual *cache.AuthEntry) {
		compareAuthEntry(t, actual, authEntry)
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigBadBase64(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(proxyEndpointScheme + proxyEndpoint),
					AuthorizationToken: aws.String(expectedUsername + ":" + expectedPassword),
				},
			},
		}, nil
	}

	credentialCache.On("Get", registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigMissingResponse(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return nil, nil
	}

	credentialCache.On("Get", registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigECRError(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return nil, nil
	}

	credentialCache.On("Get", registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigSuccessInvalidCacheHitFallback(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 1 {
			t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
		}
		return nil, nil
	}

	expiredAuthEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now().Add(-12 * time.Hour),
		ExpiresAt:          time.Now().Add(-6 * time.Hour),
		AuthorizationToken: authorizationToken,
	}

	credentialCache.On("Get", registryID).Return(expiredAuthEntry)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsSuccess(t *testing.T) {
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
	}
	authEntries := []*cache.AuthEntry{authEntry}

	credentialCache.On("List").Return(authEntries)

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, auths)
	assert.Len(t, auths, 1)

	auth := auths[0]
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsBadBase64AuthToken(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.On("List").Return(emptyCache)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 0 {
			t.Fatalf("Unexpected number of RegistryIds, expected 0 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(testProxyEndpoint),
					ExpiresAt:          aws.Time(expiresAt),
					AuthorizationToken: aws.String("invalid:token"),
				},
			},
		}, nil
	}

	auths, err := client.ListCredentials()
	assert.Error(t, err)
	assert.NotNil(t, auths)
	assert.Empty(t, auths)
}

func TestListCredentialsInvalidAuthToken(t *testing.T) {
	ecrClient := new(MockECRClient)
	credentialCache := new(cache.MockCredentialsCache)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.On("List").Return(emptyCache)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		if input == nil {
			t.Fatal("Called with nil input")
		}
		if len(input.RegistryIds) != 0 {
			t.Fatalf("Unexpected number of RegistryIds, expected 0 but got %d", len(input.RegistryIds))
		}
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []*ecr.AuthorizationData{
				{
					ProxyEndpoint:      aws.String(testProxyEndpoint),
					ExpiresAt:          aws.Time(expiresAt),
					AuthorizationToken: aws.String("invalidtoken"),
				},
			},
		}, nil
	}

	auths, err := client.ListCredentials()
	assert.Error(t, err)
	assert.NotNil(t, auths)
	assert.Empty(t, auths)
}

func compareAuthEntry(t *testing.T, actual *cache.AuthEntry, expected *cache.AuthEntry) {
	assert.NotNil(t, actual)
	assert.Equal(t, expected.AuthorizationToken, actual.AuthorizationToken)
	assert.Equal(t, expected.ProxyEndpoint, actual.ProxyEndpoint)
	assert.Equal(t, expected.ExpiresAt, actual.ExpiresAt)
	assert.WithinDuration(t, expected.RequestedAt, actual.RequestedAt, 5*time.Second)
}
