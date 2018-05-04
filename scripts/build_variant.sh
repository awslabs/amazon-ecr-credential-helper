#!/bin/bash
# Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

# This script is used for compilation of a specific variant.
# Specify GOOS as $1, GOARCH as $2
# Binaries are placed into ./bin/$GOOS-$GOARCH/docker-credential-ecr-login

# Normalize to working directory being build root (up one level from ./scripts)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
cd "${ROOT}"

# Source the shared environment
source ./scripts/shared_env

# Export variables
export TARGET_GOOS="$1"
export TARGET_GOARCH="$2"
ECR_LOGIN_VERSION="$3"
ECR_LOGIN_GITCOMMIT_SHA="$4"

./scripts/build_binary.sh "./bin/${TARGET_GOOS}-${TARGET_GOARCH}" $ECR_LOGIN_VERSION $ECR_LOGIN_GITCOMMIT_SHA

echo "Built ecr-login for ${TARGET_GOOS}-${TARGET_GOARCH}-${ECR_LOGIN_VERSION}"