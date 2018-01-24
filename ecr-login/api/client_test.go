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
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api/mocks"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	registryID       = "123456789012"
	proxyEndpoint    = "123456789012.dkr.ecr.us-east-1.amazonaws.com"
	expectedUsername = "username"
	expectedPassword = "password"
)

func TestExtractRegistry(t *testing.T) {
	testCases := []struct {
		serverURL string
		registry  *Registry
		hasError  bool
	}{
		{
			serverURL: "https://123456789012.dkr.ecr.us-east-1.amazonaws.com/v2/blah/blah",
			registry: &Registry{
				ID:     "123456789012",
				Region: "us-east-1",
			},
			hasError: false,
		},
		{
			serverURL: "123456789012.dkr.ecr.us-west-2.amazonaws.com",
			registry: &Registry{
				ID:     "123456789012",
				Region: "us-west-2",
			},
			hasError: false,
		},
		{
			serverURL: "210987654321.dkr.ecr.cn-north-1.amazonaws.com.cn/foo",
			registry: &Registry{
				ID:     "210987654321",
				Region: "cn-north-1",
			},
			hasError: false,
		},
		{
			serverURL: ".dkr.ecr.not-real.amazonaws.com",
			hasError:  true,
		},
		{
			serverURL: "not.ecr.io",
			hasError:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.serverURL, func(t *testing.T) {
			registry, err := ExtractRegistry(tc.serverURL)
			if !tc.hasError {
				assert.NoError(t, err, "No error expected")
				assert.EqualValues(t, tc.registry, registry, "Registry should be equal")
			} else {
				assert.Error(t, err, "Expected error")
			}
		})
	}
}

func TestGetAuthConfigSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			},
		},
	}, nil)

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
	}

	credentialCache.EXPECT().Get(registryID).Return(nil)
	credentialCache.EXPECT().Set(registryID, gomock.Any()).Do(
		func(_ string, actual *cache.AuthEntry) {
			compareAuthEntry(t, actual, authEntry)
		})

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigNoMatchAuthorizationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + "notproxy"),
				AuthorizationToken: aws.String(base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))),
			},
		},
	}, nil)

	credentialCache.EXPECT().Get(registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigGetCacheSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

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

	credentialCache.EXPECT().Get(registryID).Return(authEntry)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigSuccessInvalidCacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			},
		},
	}, nil)

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

	credentialCache.EXPECT().Get(registryID).Return(expiredAuthEntry)
	credentialCache.EXPECT().Set(registryID, gomock.Any()).Do(
		func(_ string, actual *cache.AuthEntry) {
			compareAuthEntry(t, actual, authEntry)
		})

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigBadBase64(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + proxyEndpoint),
				AuthorizationToken: aws.String(expectedUsername + ":" + expectedPassword),
			},
		},
	}, nil)

	credentialCache.EXPECT().Get(registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigMissingResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		})

	credentialCache.EXPECT().Get(registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigECRError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(nil, errors.New("test error"))

	credentialCache.EXPECT().Get(registryID).Return(nil)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigSuccessInvalidCacheHitFallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 1 {
				t.Fatalf("Unexpected number of RegistryIds, expected 1 but got %d", len(input.RegistryIds))
			}
		}).Return(nil, errors.New("Service eror"))

	expiredAuthEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now().Add(-12 * time.Hour),
		ExpiresAt:          time.Now().Add(-6 * time.Hour),
		AuthorizationToken: authorizationToken,
	}

	credentialCache.EXPECT().Get(registryID).Return(expiredAuthEntry)

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

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

	credentialCache.EXPECT().List().Return(authEntries)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.EXPECT().List().Return(emptyCache)

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 0 {
				t.Fatalf("Unexpected number of RegistryIds, expected 0 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String("invalid:token"),
			},
		},
	}, nil)

	auths, err := client.ListCredentials()
	assert.Error(t, err)
	assert.NotNil(t, auths)
	assert.Empty(t, auths)
}

func TestListCredentialsInvalidAuthToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecrClient := mock_ecriface.NewMockECRAPI(ctrl)
	credentialCache := mock_cache.NewMockCredentialsCache(ctrl)

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.EXPECT().List().Return(emptyCache)

	ecrClient.EXPECT().GetAuthorizationToken(gomock.Any()).Do(
		func(input *ecr.GetAuthorizationTokenInput) {
			if input == nil {
				t.Fatal("Called with nil input")
			}
			if len(input.RegistryIds) != 0 {
				t.Fatalf("Unexpected number of RegistryIds, expected 0 but got %d", len(input.RegistryIds))
			}
		}).Return(&ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String("invalidtoken"),
			},
		},
	}, nil)

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
