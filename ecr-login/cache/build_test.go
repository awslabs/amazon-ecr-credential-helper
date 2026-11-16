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

package cache

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

const (
	testRegion        = "test-region"
	testCacheFilename = "cache.json"
	testAccessKey     = "accessKey"
	testSecretKey     = "secretKey"
	testToken         = "token"
	testAccountID     = "123456789012"
	// base64 SHA-256 sum of "accountID" - FIPS-compatible
	testCredentialHash = "KjM0nn5gaorS4w48hFIfk3dFDPCQg+Fi4KmxSAzg+XI="
	// Legacy base64 MD5 sum of "accountID" for backward compatibility tests
	testLegacyCredentialHash = "MTIzNDU2Nzg5MDEy1B2M2Y8AsgTpgAmY7PhCfg=="
)

func testCredentialsProvider() aws.CredentialsProvider {
	return aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     testAccessKey,
			SecretAccessKey: testSecretKey,
			SessionToken:    testToken,
			AccountID:       testAccountID,
		}, nil
	})
}

func TestFactoryBuildFileCache(t *testing.T) {
	config := aws.Config{
		Region:      testRegion,
		Credentials: testCredentialsProvider(),
	}

	cache := BuildCredentialsCache(context.Background(), config, "")
	assert.NotNil(t, cache)

	fileCache, ok := cache.(*fileCredentialCache)

	assert.True(t, ok, "built cache is not a fileCredentialsCache")
	assert.Equal(t, fileCache.cachePrefixKey, fmt.Sprintf("%s-%s-", testRegion, testCredentialHash))
	assert.Equal(t, fileCache.filename, testCacheFilename)
}

func TestFactoryBuildNullCacheWithoutCredentials(t *testing.T) {
	config := aws.Config{
		Region:      testRegion,
		Credentials: aws.AnonymousCredentials{},
	}

	cache := BuildCredentialsCache(context.Background(), config, "")
	assert.NotNil(t, cache)

	_, ok := cache.(*nullCredentialsCache)
	assert.True(t, ok, "built cache is a nullCredentialsCache")
}

func TestFactoryBuildNullCache(t *testing.T) {
	os.Setenv("AWS_ECR_DISABLE_CACHE", "1")
	defer os.Unsetenv("AWS_ECR_DISABLE_CACHE")

	config := aws.Config{Region: testRegion}

	cache := BuildCredentialsCache(context.Background(), config, "")
	assert.NotNil(t, cache)
	_, ok := cache.(*nullCredentialsCache)
	assert.True(t, ok, "built cache is a nullCredentialsCache")
}

// TestCredentialsPrefixUsesNewHash verifies that credentialsCachePrefix uses SHA-256
func TestCredentialsPrefixUsesNewHash(t *testing.T) {
	creds := aws.Credentials{AccountID: testAccountID}
	prefix := credentialsCachePrefix(testRegion, creds)
	expectedPrefix := fmt.Sprintf("%s-%s-", testRegion, testCredentialHash)

	assert.Equal(t, expectedPrefix, prefix, "Cache prefix should use FIPS-compatible SHA-256 hash")
}

// TestIsFipsMode verifies that the isFipsMode function correctly detects FIPS mode
func TestIsFipsMode(t *testing.T) {
	tests := []struct {
		name     string
		godebug  string
		expected bool
	}{
		{"FIPS mode with fips140=on", "fips140=on", true},
		{"FIPS mode with fips140=only", "fips140=only", true},
		{"FIPS mode with other settings", "foo=bar,fips140=on,baz=qux", true},
		{"No FIPS mode", "", false},
		{"Different GODEBUG setting", "foo=bar", false},
		{"FIPS mode disabled", "fips140=off", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv("GODEBUG")
			defer func() {
				if original == "" {
					os.Unsetenv("GODEBUG")
				} else {
					os.Setenv("GODEBUG", original)
				}
			}()

			// Set test value
			if tt.godebug == "" {
				os.Unsetenv("GODEBUG")
			} else {
				os.Setenv("GODEBUG", tt.godebug)
			}

			result := isFipsMode()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLegacyKeysNotGeneratedInFipsSimulation verifies the behavior of BuildCredentialsCache
// when FIPS mode would be detected.
func TestLegacyKeysNotGeneratedInFipsSimulation(t *testing.T) {
	config := aws.Config{
		Region:      testRegion,
		Credentials: testCredentialsProvider(),
	}

	cache := BuildCredentialsCache(context.Background(), config, "")
	assert.NotNil(t, cache)

	fileCache, ok := cache.(*fileCredentialCache)
	assert.True(t, ok, "built cache should be a fileCredentialCache")

	assert.Equal(t, fmt.Sprintf("%s-%s-", testRegion, testCredentialHash), fileCache.cachePrefixKey)
	assert.Equal(t, fmt.Sprintf("%s-%s", ServiceECRPublic, testCredentialHash), fileCache.publicCacheKey)

	assert.NotEmpty(t, fileCache.legacyCachePrefixKey, "Legacy cache prefix should be present in non-FIPS mode")
	assert.NotEmpty(t, fileCache.legacyPublicCacheKey, "Legacy public cache key should be present in non-FIPS mode")
}
