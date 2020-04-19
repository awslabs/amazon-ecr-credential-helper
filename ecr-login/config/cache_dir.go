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

import "os"

func GetCacheDir() string {
	if cacheDir := os.Getenv("AWS_ECR_CACHE_DIR"); cacheDir != "" {
		return cacheDir
	}
	xdgDir := os.Getenv("XDG_CACHE_HOME");
	if xdgDir == "" {
		xdgDir = "~/.cache"
	}
	return xdgDir + "/ecr"
}
