#!/bin/sh

# Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

# A POSIX shell script which installs the required dependencies
# to build the Amazon ECR credential helper variants in a Golang
# Alpine container.

set -ex

apk add --no-cache \
    bash \
    git \
    make

# Resolves permission issues for Go cache when
# building credential helper as non-root user.
mkdir /.cache && chmod 777 /.cache
# Resolves dubious ownership of git directory when
# building credential helper as root user.
git config --global --add safe.directory /go/src/github.com/awslabs/amazon-ecr-credential-helper

