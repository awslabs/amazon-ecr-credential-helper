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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testRegistryName         = "testRegistry"
	testCachePrefixKey       = "prefix-"
	testPublicCacheKey       = "public-"
	testLegacyCachePrefixKey = "legacy-prefix-"
	testLegacyPublicCacheKey = "legacy-public-"
	testFilename             = "test.json"
)

var (
	testAuthEntry = AuthEntry{
		AuthorizationToken: "testToken",
		RequestedAt:        time.Now().Add(-5 * time.Hour),
		ExpiresAt:          time.Now().Add(7 * time.Hour),
		ProxyEndpoint:      "testEndpoint",
		Service:            ServiceECR,
	}
	testPublicAuthEntry = AuthEntry{
		AuthorizationToken: "testToken",
		RequestedAt:        time.Now().Add(-5 * time.Hour),
		ExpiresAt:          time.Now().Add(7 * time.Hour),
		ProxyEndpoint:      "testEndpoint",
		Service:            ServiceECRPublic,
	}
	testPath          = os.TempDir() + "/ecr"
	testFullFillename = filepath.Join(testPath, testFilename)
)

func TestAuthEntryValid(t *testing.T) {
	assert.True(t, testAuthEntry.IsValid(time.Now()))
}

func TestAuthEntryInValid(t *testing.T) {
	assert.True(t, testAuthEntry.IsValid(time.Now().Add(time.Second)))
}

func TestCredentials(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)

	credentialCache.Set(testRegistryName, &testAuthEntry)

	entry := credentialCache.Get(testRegistryName)
	assert.Equal(t, testAuthEntry.AuthorizationToken, entry.AuthorizationToken)
	assert.Equal(t, testAuthEntry.ProxyEndpoint, entry.ProxyEndpoint)
	assert.WithinDuration(t, testAuthEntry.RequestedAt, entry.RequestedAt, 1*time.Second)
	assert.WithinDuration(t, testAuthEntry.ExpiresAt, entry.ExpiresAt, 1*time.Second)
	assert.Equal(t, testAuthEntry.Service, entry.Service)

	entries := credentialCache.List()
	assert.NotEmpty(t, entries)
	assert.Len(t, entries, 1)
	assert.Equal(t, entry, entries[0])

	credentialCache.Clear()

	entry = credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)
}

func TestCredentialsPublic(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)

	credentialCache.Set(testRegistryName, &testPublicAuthEntry)

	entry := credentialCache.GetPublic()
	assert.Equal(t, testPublicAuthEntry.AuthorizationToken, entry.AuthorizationToken)
	assert.Equal(t, testPublicAuthEntry.ProxyEndpoint, entry.ProxyEndpoint)
	assert.WithinDuration(t, testPublicAuthEntry.RequestedAt, entry.RequestedAt, 1*time.Second)
	assert.WithinDuration(t, testPublicAuthEntry.ExpiresAt, entry.ExpiresAt, 1*time.Second)
	assert.Equal(t, testPublicAuthEntry.Service, entry.Service)

	entries := credentialCache.List()
	assert.NotEmpty(t, entries)
	assert.Len(t, entries, 1)
	assert.Equal(t, entry, entries[0])

	credentialCache.Clear()

	entry = credentialCache.GetPublic()
	assert.Nil(t, entry)
}

func TestPreviousVersionCache(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)

	registryCache := newRegistryCache()
	registryCache.Version = "0.1"
	registryCache.Registries[testRegistryName] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)

	credentialCache.Clear()
}

const testBadJson = "{nope not good json at all."

func TestInvalidCache(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)

	file, err := os.Create(testFullFillename)
	assert.NoError(t, err)

	file.WriteString(testBadJson)
	err = file.Close()
	assert.NoError(t, err)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)

	credentialCache.Clear()
}

// TestLegacyKeyBackwardCompatibility tests that credentials stored with legacy MD5-based keys
func TestLegacyKeyBackwardCompatibility(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	legacyKey := testLegacyCachePrefixKey + testRegistryName
	registryCache.Registries[legacyKey] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.NotNil(t, entry, "Should be able to retrieve credentials stored with legacy MD5 key as fallback")
	assert.Equal(t, testAuthEntry.AuthorizationToken, entry.AuthorizationToken)
	assert.Equal(t, testAuthEntry.ProxyEndpoint, entry.ProxyEndpoint)
	assert.Equal(t, testAuthEntry.Service, entry.Service)

	registryCache, err := credentialCache.(*fileCredentialCache).load()
	assert.NoError(t, err)
	assert.NotNil(t, registryCache.Registries[legacyKey], "Legacy key should still exist")
	newKey := testCachePrefixKey + testRegistryName
	assert.Nil(t, registryCache.Registries[newKey], "New key should not exist (no auto-migration)")
}

// TestLegacyPublicKeyBackwardCompatibility tests that public credentials stored with legacy MD5-based keys
func TestLegacyPublicKeyBackwardCompatibility(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	registryCache.Registries[testLegacyPublicCacheKey] = &testPublicAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.GetPublic()
	assert.NotNil(t, entry, "Should be able to retrieve public credentials stored with legacy MD5 key as fallback")
	assert.Equal(t, testPublicAuthEntry.AuthorizationToken, entry.AuthorizationToken)
	assert.Equal(t, testPublicAuthEntry.ProxyEndpoint, entry.ProxyEndpoint)
	assert.Equal(t, testPublicAuthEntry.Service, entry.Service)

	registryCache, err := credentialCache.(*fileCredentialCache).load()
	assert.NoError(t, err)
	assert.NotNil(t, registryCache.Registries[testLegacyPublicCacheKey], "Legacy public key should still exist")
	assert.Nil(t, registryCache.Registries[testPublicCacheKey], "New public key should not exist (no auto-migration)")
}

// TestNewKeyPreferredOverLegacy tests that when both new and legacy keys exist,
// the new FIPS-compatible key is preferred
func TestNewKeyPreferredOverLegacy(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	legacyKey := testLegacyCachePrefixKey + testRegistryName
	newKey := testCachePrefixKey + testRegistryName

	legacyEntry := testAuthEntry
	legacyEntry.AuthorizationToken = "legacyToken"
	registryCache.Registries[legacyKey] = &legacyEntry

	newEntry := testAuthEntry
	newEntry.AuthorizationToken = "newToken"
	registryCache.Registries[newKey] = &newEntry

	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.NotNil(t, entry)
	assert.Equal(t, "newToken", entry.AuthorizationToken, "Should prefer new FIPS-compatible key over legacy key")

	credentialCache.Clear()
}

// TestFipsModeOnlySkipsLegacyLookup tests that when GODEBUG=fips140=only is set,
// legacy MD5-based cache lookups are skipped.
func TestFipsModeOnlySkipsLegacyLookup(t *testing.T) {
	os.Setenv("GODEBUG", "fips140=only")
	defer os.Unsetenv("GODEBUG")

	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	legacyKey := testLegacyCachePrefixKey + testRegistryName
	registryCache.Registries[legacyKey] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry, "Should return nil in FIPS mode when only legacy MD5 key exists")
}

// TestFipsModeOnSkipsLegacyLookup tests that when GODEBUG=fips140=on is set,
// legacy MD5-based cache lookups are skipped.
func TestFipsModeOnSkipsLegacyLookup(t *testing.T) {
	os.Setenv("GODEBUG", "fips140=on")
	defer os.Unsetenv("GODEBUG")

	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	legacyKey := testLegacyCachePrefixKey + testRegistryName
	registryCache.Registries[legacyKey] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	// Try to retrieve - should return nil because FIPS mode skips MD5 lookup
	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry, "Should return nil in FIPS mode when only legacy MD5 key exists")
}

// TestFipsModeOnlySkipsLegacyPublicLookup tests that when GODEBUG=fips140=only is set,
// legacy MD5-based public cache lookups are skipped.
func TestFipsModeOnlySkipsLegacyPublicLookup(t *testing.T) {
	os.Setenv("GODEBUG", "fips140=only")
	defer os.Unsetenv("GODEBUG")

	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	registryCache.Registries[testLegacyPublicCacheKey] = &testPublicAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.GetPublic()
	assert.Nil(t, entry, "Should return nil in FIPS mode when only legacy MD5 key exists for public")
}

// TestFipsModeWithNewKey tests that FIPS mode still works with SHA-256 keys
func TestFipsModeWithNewKey(t *testing.T) {
	os.Setenv("GODEBUG", "fips140=only")
	defer os.Unsetenv("GODEBUG")

	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey, testLegacyCachePrefixKey, testLegacyPublicCacheKey)
	defer credentialCache.Clear()

	registryCache := newRegistryCache()
	newKey := testCachePrefixKey + testRegistryName
	registryCache.Registries[newKey] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.NotNil(t, entry, "Should find credentials with SHA-256 key in FIPS mode")
	assert.Equal(t, testAuthEntry.AuthorizationToken, entry.AuthorizationToken)
}
