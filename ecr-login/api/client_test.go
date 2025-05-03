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
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	ecrpublictypes "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	mock_api "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api/mocks"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	mock_cache "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache/mocks"
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
	}{{
		serverURL: "https://123456789012.dkr.ecr.us-east-1.amazonaws.com/v2/blah/blah",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    false,
			Region:  "us-east-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "123456789012.dkr.ecr.us-west-2.amazonaws.com",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    false,
			Region:  "us-west-2",
			Service: ServiceECR,
		},
		hasError: false,
        }, {
                serverURL: "123456789012.dkr-ecr.us-west-2.on.aws",
                registry: &Registry{
                        ID:      "123456789012",
                        FIPS:    false,
                        Region:  "us-west-2",
                        Service: ServiceECR,
                },
                hasError: false,
	}, {
		serverURL: "210987654321.dkr.ecr.cn-north-1.amazonaws.com.cn/foo",
		registry: &Registry{
			ID:      "210987654321",
			FIPS:    false,
			Region:  "cn-north-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "210987654321.dkr.ecr.us-iso-east-1.c2s.ic.gov",
		registry: &Registry{
			ID:      "210987654321",
			FIPS:    false,
			Region:  "us-iso-east-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "123456789012.dkr.ecr.us-isob-east-1.sc2s.sgov.gov",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    false,
			Region:  "us-isob-east-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "123456789012.dkr.ecr.eu-isoe-west-1.cloud.adc-e.uk",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    false,
			Region:  "eu-isoe-west-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "123456789012.dkr.ecr.us-isof-east-1.csp.hci.ic.gov",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    false,
			Region:  "us-isof-east-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "123456789012.dkr.ecr-fips.us-gov-west-1.amazonaws.com",
		registry: &Registry{
			ID:      "123456789012",
			FIPS:    true,
			Region:  "us-gov-west-1",
			Service: ServiceECR,
		},
		hasError: false,
	}, {
		serverURL: "https://public.ecr.aws",
		registry: &Registry{
			Service: ServiceECRPublic,
		},
	}, {
		serverURL: "public.ecr.aws",
		registry: &Registry{
			Service: ServiceECRPublic,
		},
	}, {
		serverURL: "https://public.ecr.aws/amazonlinux",
		registry: &Registry{
			Service: ServiceECRPublic,
		},
	}, {
		serverURL: ".dkr.ecr.not-real.amazonaws.com",
		hasError:  true,
	}, {
		serverURL: "not.ecr.io",
		hasError:  true,
	}, {
		serverURL: "https://123456789012.dkr.ecr.us-west-2.amazonaws.com.fake.example.com/image:latest",
		hasError:  true,
	}, {
		serverURL: "123456789012.dkr.ecr.us-west-2.amazonaws.com.fake.example.com",
		hasError:  true,
	}, {
		serverURL: "123456789012.dkr.ecr-fips.us-gov-west-1.amazonaws.com.fake.example.com",
		hasError:  true,
	}, {
		serverURL: "210987654321.dkr.ecr.cn-north-1.amazonaws.com.cn.fake.example.com.cn",
		hasError:  true,
	}, {
		serverURL: "https://public.ecr.aws.fake.example.com",
		hasError:  true,
	}, {
		serverURL: "public.ecr.aws.fake.example.com",
		hasError:  true,
	}}
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
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			}},
		}, nil
	}

	authEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
	}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }
	credentialCache.SetFn = func(_ string, actual *cache.AuthEntry) {
		compareAuthEntry(t, actual, authEntry)
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigNoMatchAuthorizationToken(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + "notproxy"),
				AuthorizationToken: aws.String(base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))),
			}},
		}, nil
	}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigGetCacheSuccess(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

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
	credentialCache.GetFn = func(r string) *cache.AuthEntry {
		assert.Equal(t, registryID, r, "get from cache")
		return authEntry
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigSuccessInvalidCacheHit(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			}},
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

	credentialCache.GetFn = func(r string) *cache.AuthEntry {
		assert.Equal(t, registryID, r, "get from cache")
		return expiredAuthEntry
	}
	credentialCache.SetFn = func(_ string, actual *cache.AuthEntry) {
		compareAuthEntry(t, actual, authEntry)
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestGetAuthConfigBadBase64(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(proxyEndpointScheme + proxyEndpoint),
				AuthorizationToken: aws.String(expectedUsername + ":" + expectedPassword),
			}},
		}, nil
	}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigMissingResponse(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return nil, nil
	}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigECRError(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return nil, errors.New("test error")
	}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Nil(t, auth)
}

func TestGetAuthConfigSuccessInvalidCacheHitFallback(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		ecrClient:       ecrClient,
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 1, "GetAuthorizationToken registry IDs len")
		return nil, errors.New("service error")
	}

	expiredAuthEntry := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now().Add(-12 * time.Hour),
		ExpiresAt:          time.Now().Add(-6 * time.Hour),
		AuthorizationToken: authorizationToken,
	}

	credentialCache.GetFn = func(r string) *cache.AuthEntry {
		assert.Equal(t, registryID, r, "get from cache")
		return expiredAuthEntry
	}

	auth, err := client.GetCredentials(proxyEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsSuccess(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	ecrPublicClient := &mock_api.MockECRPublicAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		credentialCache: credentialCache,
		ecrClient:       ecrClient,
		ecrPublicClient: ecrPublicClient,
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

	ecrClient.GetAuthorizationTokenFn = func(_ *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		return nil, errors.New("test error")
	}
	ecrPublicClient.GetAuthorizationTokenFn = func(_ *ecrpublic.GetAuthorizationTokenInput) (*ecrpublic.GetAuthorizationTokenOutput, error) {
		return nil, errors.New("test error")
	}
	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }
	credentialCache.GetPublicFn = func() *cache.AuthEntry { return nil }
	credentialCache.ListFn = func() []*cache.AuthEntry { return authEntries }

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, auths)
	assert.Len(t, auths, 1)

	auth := auths[0]
	assert.Equal(t, auth.Username, expectedUsername)
	assert.Equal(t, auth.Password, expectedPassword)
	assert.Equal(t, auth.ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsCached(t *testing.T) {
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		credentialCache: credentialCache,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	authEntry1 := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
		Service:            cache.ServiceECR,
	}
	authEntry2 := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
		Service:            cache.ServiceECRPublic,
	}
	authEntries := []*cache.AuthEntry{authEntry1, authEntry2}

	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return authEntry1 }
	credentialCache.GetPublicFn = func() *cache.AuthEntry { return authEntry2 }
	credentialCache.ListFn = func() []*cache.AuthEntry { return authEntries }

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, auths)
	assert.Len(t, auths, 2)

	assert.Equal(t, auths[0].Username, expectedUsername)
	assert.Equal(t, auths[0].Password, expectedPassword)
	assert.Equal(t, auths[0].ProxyEndpoint, testProxyEndpoint)
	assert.Equal(t, auths[1].Username, expectedUsername)
	assert.Equal(t, auths[1].Password, expectedPassword)
	assert.Equal(t, auths[1].ProxyEndpoint, testProxyEndpoint)
}

func TestListCredentialsEmpty(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	ecrPublicClient := &mock_api.MockECRPublicAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		credentialCache: credentialCache,
		ecrClient:       ecrClient,
		ecrPublicClient: ecrPublicClient,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	authorizationToken := base64.StdEncoding.EncodeToString([]byte(expectedUsername + ":" + expectedPassword))
	expiresAt := time.Now().Add(12 * time.Hour)

	authEntry1 := &cache.AuthEntry{
		ProxyEndpoint:      testProxyEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
		Service:            cache.ServiceECR,
	}
	authEntry2 := &cache.AuthEntry{
		ProxyEndpoint:      ecrPublicEndpoint,
		RequestedAt:        time.Now(),
		ExpiresAt:          expiresAt,
		AuthorizationToken: authorizationToken,
		Service:            cache.ServiceECRPublic,
	}
	authEntries := []*cache.AuthEntry{authEntry1, authEntry2}

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 0, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			}},
		}, nil
	}
	ecrPublicClient.GetAuthorizationTokenFn = func(*ecrpublic.GetAuthorizationTokenInput) (*ecrpublic.GetAuthorizationTokenOutput, error) {
		return &ecrpublic.GetAuthorizationTokenOutput{
			AuthorizationData: &ecrpublictypes.AuthorizationData{
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String(authorizationToken),
			},
		}, nil
	}
	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }
	credentialCache.GetPublicFn = func() *cache.AuthEntry { return nil }
	credentialCache.ListFn = func() []*cache.AuthEntry { return authEntries }
	setCallCount := 0
	credentialCache.SetFn = func(_ string, _ *cache.AuthEntry) {
		setCallCount++
	}

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, auths)
	assert.Len(t, auths, 2)
	assert.Equal(t, 2, setCallCount)

	assert.Equal(t, auths[0].Username, expectedUsername)
	assert.Equal(t, auths[0].Password, expectedPassword)
	assert.Equal(t, auths[0].ProxyEndpoint, testProxyEndpoint)
	assert.Equal(t, auths[1].Username, expectedUsername)
	assert.Equal(t, auths[1].Password, expectedPassword)
	assert.Equal(t, auths[1].ProxyEndpoint, ecrPublicEndpoint)
}

func TestListCredentialsBadBase64AuthToken(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	ecrPublicClient := &mock_api.MockECRPublicAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		credentialCache: credentialCache,
		ecrClient:       ecrClient,
		ecrPublicClient: ecrPublicClient,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.ListFn = func() []*cache.AuthEntry { return emptyCache }
	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }
	credentialCache.GetPublicFn = func() *cache.AuthEntry { return nil }

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 0, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String("invalid:token"),
			}},
		}, nil
	}
	ecrPublicClient.GetAuthorizationTokenFn = func(_ *ecrpublic.GetAuthorizationTokenInput) (*ecrpublic.GetAuthorizationTokenOutput, error) {
		return nil, errors.New("test error")
	}

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, auths)
	assert.Empty(t, auths)
}

func TestListCredentialsInvalidAuthToken(t *testing.T) {
	ecrClient := &mock_api.MockECRAPI{}
	ecrPublicClient := &mock_api.MockECRPublicAPI{}
	credentialCache := &mock_cache.MockCredentialsCache{}

	client := &defaultClient{
		credentialCache: credentialCache,
		ecrClient:       ecrClient,
		ecrPublicClient: ecrPublicClient,
	}

	testProxyEndpoint := proxyEndpointScheme + proxyEndpoint
	expiresAt := time.Now().Add(12 * time.Hour)

	emptyCache := []*cache.AuthEntry{}
	credentialCache.ListFn = func() []*cache.AuthEntry { return emptyCache }
	credentialCache.GetFn = func(_ string) *cache.AuthEntry { return nil }
	credentialCache.GetPublicFn = func() *cache.AuthEntry { return nil }

	ecrClient.GetAuthorizationTokenFn = func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
		assert.NotNil(t, input, "GetAuthorizationToken input")
		assert.Len(t, input.RegistryIds, 0, "GetAuthorizationToken registry IDs len")
		return &ecr.GetAuthorizationTokenOutput{
			AuthorizationData: []ecrtypes.AuthorizationData{{
				ProxyEndpoint:      aws.String(testProxyEndpoint),
				ExpiresAt:          aws.Time(expiresAt),
				AuthorizationToken: aws.String("invalidtoken"),
			}},
		}, nil
	}
	ecrPublicClient.GetAuthorizationTokenFn = func(_ *ecrpublic.GetAuthorizationTokenInput) (*ecrpublic.GetAuthorizationTokenOutput, error) {
		return nil, errors.New("test error")
	}

	auths, err := client.ListCredentials()
	assert.NoError(t, err)
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
