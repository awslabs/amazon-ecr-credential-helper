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

package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetCacheDir returns the cache directory to use for the credential helper.
// The cache directory is determined in this order:
// 1. The location specified by the AWS_ECR_CACHE_DIR environment variable, if
//    set and if expanded validly by homedir.Expand.
// 2. ~/.ecr, if it already exists and is a directory, for
//    backwards-compatibility.
// 3. The /ecr directory under the value returned by os.UserCacheDir as the
//    default value. os.UserCacheDir returns the OS-specific cache directory.
func GetCacheDir() (string, error) {
	var (
		cacheDir string
		userCacheDir string
		err error
	)

	if cacheDir = os.Getenv("AWS_ECR_CACHE_DIR"); cacheDir != "" {
		if cacheDir, err = homedir.Expand(cacheDir); err == nil {
			return cacheDir, nil
		}
	}

	cacheDir = "~/.ecr"
	if cacheDir, err = homedir.Expand(cacheDir); err == nil {
		if info, err := os.Stat(cacheDir); err == nil {
			if info.IsDir() {
				return cacheDir, nil
			}
		}
	}

	userCacheDir, err = os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userCacheDir,"/ecr"), nil
}
