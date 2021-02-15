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

const (
	testRegistryName   = "testRegistry"
	testCachePrefixKey = "prefix-"
	testPublicCacheKey = "public-"
	testFilename       = "test.json"
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
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey)

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
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey)

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
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey)

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
	credentialCache := NewFileCredentialsCache(testPath, testFilename, testCachePrefixKey, testPublicCacheKey)

	file, err := os.Create(testFullFillename)
	assert.NoError(t, err)

	file.WriteString(testBadJson)
	err = file.Close()
	assert.NoError(t, err)

	entry := credentialCache.Get(testRegistryName)
	assert.Nil(t, entry)

	credentialCache.Clear()
}
