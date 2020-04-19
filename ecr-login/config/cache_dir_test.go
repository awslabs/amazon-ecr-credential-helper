// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCacheDir(t *testing.T) {
	// caching is disabled so we can adjust the HOME env var
	homedir.DisableCache = true
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)
	restore := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", restore)
	home, err := homedir.Dir()
	require.NoError(t, err, "home must be resolved")
	require.Equal(t, tempDir, home, "expect home to be tempDir")
	userCache, err := os.UserCacheDir()
	require.NoError(t, err, "user cache dir must be resolved")
	// Note: tests that adjust environment variables are not thread-safe
	t.Run("AWS_ECR_CACHE_DIR valid", func(t *testing.T) {
		restore := os.Getenv("AWS_ECR_CACHE_DIR")
		os.Setenv("AWS_ECR_CACHE_DIR", "~/testval")
		defer os.Setenv("AWS_ECR_CACHE_DIR", restore)

		actual, err := GetCacheDir()
		assert.NoError(t, err, "should resolve cache dir")
		assert.Equal(t, filepath.Join(home, "testval"), actual)
	})

	t.Run("AWS_ECR_CACHE_DIR invalid", func(t *testing.T) {
		restore := os.Getenv("AWS_ECR_CACHE_DIR")
		os.Setenv("AWS_ECR_CACHE_DIR", "~testval")
		defer os.Setenv("AWS_ECR_CACHE_DIR", restore)

		actual, err := GetCacheDir()
		assert.NoError(t, err, "should resolve cache dir")
		assert.Equal(t, filepath.Join(userCache, "ecr"), actual)
	})

	t.Run("old exists", func(t *testing.T) {
		err := os.Mkdir(filepath.Join(tempDir, ".ecr"), 0755)
		require.NoError(t, err)
		defer os.RemoveAll(filepath.Join(tempDir, ".ecr"))

		actual, err := GetCacheDir()
		assert.NoError(t, err, "should resolve cache dir")
		assert.Equal(t, filepath.Join(tempDir, ".ecr"), actual)
	})

	t.Run("default", func(t *testing.T) {
		actual, err := GetCacheDir()
		assert.NoError(t, err, "should resolve cache dir")
		assert.Equal(t, filepath.Join(userCache, "ecr"), actual)
	})
}