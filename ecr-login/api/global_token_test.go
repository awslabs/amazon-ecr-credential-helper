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
	"encoding/base64"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	ecrpublictypes "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/aws/aws-sdk-go-v2/service/signin"
	mock_api "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api/mocks"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
	mock_cache "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache/mocks"
	ecrconfig "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/config"
	"github.com/stretchr/testify/assert"
)

func TestGetCredentialsUsesGlobalTokenWhenSelected(t *testing.T) {
	t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", ecrconfig.TokenTypeGlobal)
	t.Setenv("AWS_USE_FIPS_ENDPOINT", "")
	const registryHostname = "123456789012.dkr.ecr.us-west-2.amazonaws.com"
	const rawAccessToken = "raw-signin-access-token"

	type contextKey string
	const requestContextKey contextKey = "request"
	ctx := context.WithValue(context.Background(), requestContextKey, "request-value")

	credentialCache := &mock_cache.MockCredentialsCache{
		GetFn: func(key string) *cache.AuthEntry {
			assert.Equal(t, registryHostname, key)
			return nil
		},
		SetFn: func(key string, entry *cache.AuthEntry) {
			assert.Equal(t, registryHostname, key)
			assert.Equal(t, cache.ServiceECR, entry.Service)
			assert.Equal(t, "https://"+registryHostname, entry.ProxyEndpoint)
			assert.WithinDuration(t, entry.RequestedAt.Add(time.Hour), entry.ExpiresAt, time.Second)
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("AWS:"+rawAccessToken)), entry.AuthorizationToken)
		},
	}
	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(callCtx context.Context, input *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			assert.Equal(t, "request-value", callCtx.Value(requestContextKey))
			assert.Equal(t, signInGrantType, aws.ToString(input.GrantType))
			assert.Equal(t, signInResource, aws.ToString(input.Resource))
			assert.Nil(t, input.DurationSeconds)
			return &signin.CreateBearerTokenOutput{
				AccessToken: aws.String(rawAccessToken),
				ExpiresIn:   aws.Int32(3600),
			}, nil
		},
	}
	client := &defaultClient{
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	auth, err := client.GetCredentials(ctx, registryHostname)

	assert.NoError(t, err)
	assert.Equal(t, "AWS", auth.Username)
	assert.Equal(t, rawAccessToken, auth.Password)
	assert.Equal(t, "https://"+registryHostname, auth.ProxyEndpoint)
}

func TestGetGlobalCredentialsUsesValidCache(t *testing.T) {
	const registryHostname = "123456789012.dkr.ecr.us-west-2.amazonaws.com"
	cachedToken := base64.StdEncoding.EncodeToString([]byte("AWS:cached-global-token"))
	credentialCache := &mock_cache.MockCredentialsCache{
		GetFn: func(key string) *cache.AuthEntry {
			assert.Equal(t, registryHostname, key)
			return &cache.AuthEntry{
				AuthorizationToken: cachedToken,
				RequestedAt:        time.Now(),
				ExpiresAt:          time.Now().Add(time.Hour),
				ProxyEndpoint:      "https://" + registryHostname,
				Service:            cache.ServiceECR,
			}
		},
		SetFn: func(_ string, _ *cache.AuthEntry) {
			t.Fatal("a valid global token must not be replaced")
		},
	}
	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			t.Fatal("SignIn must not be called for a valid cached global token")
			return nil, nil
		},
	}
	client := &defaultClient{
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	auth, err := client.GetGlobalCredentials(context.Background(), registryHostname, "us-west-2")

	assert.NoError(t, err)
	assert.Equal(t, "cached-global-token", auth.Password)
	assert.Equal(t, "https://"+registryHostname, auth.ProxyEndpoint)
}

func TestGetGlobalCredentialsFallsBackToExpiredCache(t *testing.T) {
	const registryHostname = "123456789012.dkr.ecr.us-west-2.amazonaws.com"
	cachedToken := base64.StdEncoding.EncodeToString([]byte("AWS:cached-global-token"))
	credentialCache := &mock_cache.MockCredentialsCache{
		GetFn: func(key string) *cache.AuthEntry {
			assert.Equal(t, registryHostname, key)
			return &cache.AuthEntry{
				AuthorizationToken: cachedToken,
				RequestedAt:        time.Now().Add(-2 * time.Hour),
				ExpiresAt:          time.Now().Add(-time.Hour),
				ProxyEndpoint:      "https://" + registryHostname,
				Service:            cache.ServiceECR,
			}
		},
		SetFn: func(_ string, _ *cache.AuthEntry) {
			t.Fatal("a failed token request must not replace the cached token")
		},
	}
	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			return nil, errors.New("service unavailable")
		},
	}
	client := &defaultClient{
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	auth, err := client.GetGlobalCredentials(context.Background(), registryHostname, "us-west-2")

	assert.NoError(t, err)
	assert.Equal(t, "cached-global-token", auth.Password)
}

func TestGetGlobalCredentialsRejectsInvalidResponses(t *testing.T) {
	testCases := []struct {
		name     string
		response *signin.CreateBearerTokenOutput
		apiError error
		error    string
	}{
		{name: "API error", apiError: errors.New("access denied"), error: "access denied"},
		{name: "nil response", error: "missing or invalid token data"},
		{
			name: "empty access token",
			response: &signin.CreateBearerTokenOutput{
				ExpiresIn: aws.Int32(3600),
			},
			error: "missing or invalid token data",
		},
		{
			name: "invalid expiration",
			response: &signin.CreateBearerTokenOutput{
				AccessToken: aws.String("test-access-token"),
			},
			error: "missing or invalid token data",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			credentialCache := &mock_cache.MockCredentialsCache{
				GetFn: func(_ string) *cache.AuthEntry {
					return nil
				},
				SetFn: func(_ string, _ *cache.AuthEntry) {
					t.Fatal("an invalid response must not be cached")
				},
			}
			signInClient := &mock_api.MockSignInAPI{
				CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
					return testCase.response, testCase.apiError
				},
			}
			client := &defaultClient{
				signInClient:    signInClient,
				credentialCache: credentialCache,
			}

			auth, err := client.GetGlobalCredentials(context.Background(), "123456789012.dkr.ecr.us-west-2.amazonaws.com", "us-west-2")

			assert.Nil(t, auth)
			assert.ErrorContains(t, err, testCase.error)
		})
	}
}

func TestGlobalTokenCacheIsScopedByHostname(t *testing.T) {
	requestedHostnames := []string{
		"111111111111.dkr.ecr.us-west-2.amazonaws.com",
		"222222222222.dkr.ecr.us-west-2.amazonaws.com",
	}
	getKeys := make([]string, 0, len(requestedHostnames))
	setKeys := make([]string, 0, len(requestedHostnames))
	credentialCache := &mock_cache.MockCredentialsCache{
		GetFn: func(key string) *cache.AuthEntry {
			getKeys = append(getKeys, key)
			return nil
		},
		SetFn: func(key string, _ *cache.AuthEntry) {
			setKeys = append(setKeys, key)
		},
	}
	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			return &signin.CreateBearerTokenOutput{
				AccessToken: aws.String("test-access-token"),
				ExpiresIn:   aws.Int32(3600),
			}, nil
		},
	}
	client := &defaultClient{
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	for _, hostname := range requestedHostnames {
		_, err := client.GetGlobalCredentials(context.Background(), hostname, "us-west-2")
		assert.NoError(t, err)
	}

	assert.Equal(t, requestedHostnames, getKeys)
	assert.Equal(t, requestedHostnames, setKeys)
}

func TestCanonicalProxyEndpoint(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hostname",
			input:    "123456789012.dkr.ecr.us-west-2.amazonaws.com",
			expected: "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
		},
		{
			name:     "scheme path and port",
			input:    "https://123456789012.dkr.ecr.us-west-2.amazonaws.com:443/repository:tag",
			expected: "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			endpoint, err := canonicalProxyEndpoint(testCase.input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, endpoint)
		})
	}
}

func TestGetCredentialsRejectsInvalidTokenType(t *testing.T) {
	t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", "signin")

	client := &defaultClient{}
	auth, err := client.GetCredentials(context.Background(), "123456789012.dkr.ecr.us-west-2.amazonaws.com")

	assert.Nil(t, auth)
	assert.ErrorContains(t, err, "unsupported AWS_ECR_PRIVATE_TOKEN_TYPE")
}

func TestTokenTypeValidationDoesNotApplyToECRPublic(t *testing.T) {
	t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", "unsupported")

	authorizationToken := base64.StdEncoding.EncodeToString([]byte("AWS:public-token"))
	credentialCache := &mock_cache.MockCredentialsCache{
		GetPublicFn: func() *cache.AuthEntry {
			return nil
		},
		SetFn: func(key string, entry *cache.AuthEntry) {
			assert.Equal(t, ecrPublicName, key)
			assert.Equal(t, cache.ServiceECRPublic, entry.Service)
		},
	}
	ecrPublicClient := &mock_api.MockECRPublicAPI{
		GetAuthorizationTokenFn: func(_ *ecrpublic.GetAuthorizationTokenInput) (*ecrpublic.GetAuthorizationTokenOutput, error) {
			return &ecrpublic.GetAuthorizationTokenOutput{
				AuthorizationData: &ecrpublictypes.AuthorizationData{
					AuthorizationToken: aws.String(authorizationToken),
					ExpiresAt:          aws.Time(time.Now().Add(12 * time.Hour)),
				},
			}, nil
		},
	}
	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			t.Fatal("ECR Public must not call SignIn")
			return nil, nil
		},
	}
	client := &defaultClient{
		ecrPublicClient: ecrPublicClient,
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	auth, err := client.GetCredentials(context.Background(), ecrPublicName)

	assert.NoError(t, err)
	assert.Equal(t, "public-token", auth.Password)
}

func TestInvalidTokenTypeDoesNotChangeUnknownHost(t *testing.T) {
	t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", "unsupported")

	client := &defaultClient{}
	auth, err := client.GetCredentials(context.Background(), "registry.example.com")

	assert.Nil(t, auth)
	assert.ErrorContains(t, err, "can only be used with Amazon Elastic Container Registry")
}

func TestGlobalTokenFIPSInputsDoNotBlockSignInCall(t *testing.T) {
	testCases := []struct {
		name       string
		serverURL  string
		useFIPSEnv string
	}{
		{
			name:      "FIPS registry hostname",
			serverURL: "123456789012.dkr.ecr-fips.us-gov-west-1.amazonaws.com",
		},
		{
			name:       "FIPS environment setting",
			serverURL:  "123456789012.dkr.ecr.us-west-2.amazonaws.com",
			useFIPSEnv: "true",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", ecrconfig.TokenTypeGlobal)
			t.Setenv("AWS_USE_FIPS_ENDPOINT", testCase.useFIPSEnv)

			calls := 0
			signInClient := &mock_api.MockSignInAPI{
				CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
					calls++
					return &signin.CreateBearerTokenOutput{
						AccessToken: aws.String("test-access-token"),
						ExpiresIn:   aws.Int32(3600),
					}, nil
				},
			}
			credentialCache := &mock_cache.MockCredentialsCache{
				GetFn: func(_ string) *cache.AuthEntry {
					return nil
				},
				SetFn: func(_ string, _ *cache.AuthEntry) {},
			}
			client := &defaultClient{
				signInClient:    signInClient,
				credentialCache: credentialCache,
			}

			auth, err := client.GetCredentials(context.Background(), testCase.serverURL)

			assert.NoError(t, err)
			assert.Equal(t, "test-access-token", auth.Password)
			assert.Equal(t, 1, calls)
		})
	}
}

func TestGlobalTokenSignInErrorIncludesRegistryAndRegion(t *testing.T) {
	t.Setenv("AWS_ECR_PRIVATE_TOKEN_TYPE", ecrconfig.TokenTypeGlobal)
	t.Setenv("AWS_USE_FIPS_ENDPOINT", "")
	const registryHostname = "123456789012.dkr.ecr.us-west-2.amazonaws.com"
	signInError := errors.New("SDK endpoint resolution failed")

	signInClient := &mock_api.MockSignInAPI{
		CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
			return nil, signInError
		},
	}
	credentialCache := &mock_cache.MockCredentialsCache{
		GetFn: func(_ string) *cache.AuthEntry {
			return nil
		},
	}
	client := &defaultClient{
		signInClient:    signInClient,
		credentialCache: credentialCache,
	}

	auth, err := client.GetCredentials(context.Background(), registryHostname)

	assert.Nil(t, auth)
	assert.ErrorContains(t, err, `signin: failed to create bearer token for registry "`+registryHostname+`" (region "us-west-2")`)
	assert.ErrorIs(t, err, signInError)
}

func TestDefaultClientFactoryConfiguresSignInRegion(t *testing.T) {
	t.Setenv("AWS_ECR_DISABLE_CACHE", "true")

	client, err := DefaultClientFactory{}.NewClient(context.Background(), aws.Config{
		Region:      "us-west-2",
		Credentials: aws.AnonymousCredentials{},
	})
	assert.NoError(t, err)
	defaultClient, ok := client.(*defaultClient)
	if !assert.True(t, ok) {
		return
	}
	signInWrapper, ok := defaultClient.signInClient.(*signInClientWrapper)
	if !assert.True(t, ok) {
		return
	}
	signInClient, ok := signInWrapper.client.(*signin.Client)
	if !assert.True(t, ok) {
		return
	}

	assert.Equal(t, "us-west-2", signInClient.Options().Region)
}

func TestSignInClientWrapper(t *testing.T) {
	const sensitiveValue = "test-query-value"
	regularErr := errors.New("service unavailable")
	testCases := []struct {
		name     string
		apiError error
		verify   func(*testing.T, error)
	}{
		{
			name: "sanitizes url.Error",
			apiError: &url.Error{
				Op:  "Post",
				URL: "https://signin.us-west-2.amazonaws.com/?request=" + sensitiveValue,
				Err: errors.New("connection refused"),
			},
			verify: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.NotContains(t, err.Error(), sensitiveValue)
				assert.Contains(t, err.Error(), "request=redacted")
				assert.Contains(t, err.Error(), "connection refused")

				var resultURLErr *url.Error
				assert.True(t, errors.As(err, &resultURLErr))
			},
		},
		{
			name:     "passes through non-url.Error unchanged",
			apiError: regularErr,
			verify: func(t *testing.T, err error) {
				assert.Equal(t, regularErr, err)
			},
		},
		{
			name: "returns nil error when no error",
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mock := &mock_api.MockSignInAPI{
				CreateBearerTokenFn: func(_ context.Context, _ *signin.CreateBearerTokenInput) (*signin.CreateBearerTokenOutput, error) {
					return nil, testCase.apiError
				},
			}
			client := NewSignInClientWrapper(mock)

			_, err := client.CreateBearerToken(context.Background(), &signin.CreateBearerTokenInput{})

			testCase.verify(t, err)
		})
	}
}
