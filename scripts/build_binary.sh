#!/bin/bash
# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
# 	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

# Normalize to working directory being build root (up one level from ./scripts)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
cd "${ROOT}"

# Builds the ecr-login binary from source in the specified destination paths.
mkdir -p $1

cd "${ROOT}"

package_root="github.com/awslabs/amazon-ecr-credential-helper/ecr-login"

version_ldflags=""
if [[ -n "${2}" ]]; then
  version_ldflags="-X ${package_root}/version.Version=${2}"
fi

if [[ -n "${3}" ]]; then
  version_ldflags="$version_ldflags -X ${package_root}/version.GitCommitSHA=${3}"
fi

GOOS=$TARGET_GOOS GOARCH=$TARGET_GOARCH CGO_ENABLED=0 \
       	go build -installsuffix cgo -a -ldflags "-s ${version_ldflags}" \
       	-o $1/docker-credential-ecr-login \
	./ecr-login/cli/docker-credential-ecr-login
