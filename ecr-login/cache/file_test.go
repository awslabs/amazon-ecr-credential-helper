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

package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testAuthEntry = AuthEntry{
	AuthorizationToken: "testToken",
	RequestedAt:        time.Now().Add(-5 * time.Hour),
	ExpiresAt:          time.Now().Add(7 * time.Hour),
	ProxyEndpoint:      "testEndpoint",
}

var testRegistryName = "testRegistry"

var testCachePrefixKey = "prefix-"
var testPath = os.TempDir() + "/ecr"
var testFilename = "test.json"
var testFullFilename = filepath.Join(testPath, testFilename)

func TestAuthEntryValid(t *testing.T) {
	assert.True(t, testAuthEntry.IsValid(time.Now()))
}

func TestAuthEntryInValid(t *testing.T) {
	assert.True(t, testAuthEntry.IsValid(time.Now().Add(time.Second)))
}

// TestCredentials tests the credentials was successfully save and loaded
func TestCredentials(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey)

	credentialCache.Set(testRegistryName, &testAuthEntry)

	entry := credentialCache.Get(testRegistryName)
	assert.NotNil(t, entry)
	assert.Equal(t, testAuthEntry.AuthorizationToken, entry.AuthorizationToken)
	assert.Equal(t, testAuthEntry.ProxyEndpoint, entry.ProxyEndpoint)
	assert.WithinDuration(t, testAuthEntry.RequestedAt, entry.RequestedAt, 1*time.Second)
	assert.WithinDuration(t, testAuthEntry.ExpiresAt, entry.ExpiresAt, 1*time.Second)

	entries := credentialCache.List()
	assert.NotEmpty(t, entries)
	assert.Len(t, entries, 1)
	assert.Equal(t, entry, entries[0])

	credentialCache.Clear()

	entry = credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)
}

// TestPreviousVersionCache tests the previous version of cache was not loaded
// with current version of registry cache
func TestPreviousVersionCache(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey)

	registryCache := newRegistryCache()
	registryCache.Version = "0.1"
	registryCache.Registries[testRegistryName] = &testAuthEntry
	credentialCache.(*fileCredentialCache).save(registryCache)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)

	credentialCache.Clear()
}

const testBadJSON = "{nope not good json at all."

// TestInvalidCache tests no credentials will be loaded from invalid cache file
func TestInvalidCache(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey)

	file, err := os.Create(testFullFilename)
	assert.NoError(t, err)

	file.WriteString(testBadJSON)
	err = file.Close()
	assert.NoError(t, err)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)

	credentialCache.Clear()
}

// TestCleanupExpiredCredentialsOnSave tests the expired credentials was not cached in file
func TestCleanupExpiredCredentialsOnSave(t *testing.T) {
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey)

	registryCache := newRegistryCache()

	// The second auth entry is expired by the time
	testAuthEntry2 := testAuthEntry
	testAuthEntry2.ExpiresAt = time.Now().Add(-1 * time.Second)
	registryCache.Registries["testRegistry1"] = &testAuthEntry
	registryCache.Registries["testRegistry2"] = &testAuthEntry2

	credentialCache.(*fileCredentialCache).save(registryCache)
	defer credentialCache.Clear()

	entries := credentialCache.List()
	assert.Len(t, entries, 1)
}

// TestCleanupExpiredCredentialsOnSave tests the expired credentials was cached
// when the clean up expired credentials was disabled
func TestCleanupExpiredCredentialsDisabled(t *testing.T) {
	os.Setenv(cleanExpiredCredentialsEnv, "true")
	defer os.Unsetenv(cleanExpiredCredentialsEnv)
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey)

	registryCache := newRegistryCache()

	// The second auth entry is expired by the time
	testAuthEntry2 := testAuthEntry
	testAuthEntry2.ExpiresAt = time.Now().Add(-1 * time.Second)
	registryCache.Registries["testRegistry1"] = &testAuthEntry
	registryCache.Registries["testRegistry2"] = &testAuthEntry2

	credentialCache.(*fileCredentialCache).save(registryCache)
	defer credentialCache.Clear()

	entries := credentialCache.List()
	assert.Len(t, entries, 2)
}
