#!/bin/bash
# Copyright 2017-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

# This script checks out a git repository with the correct revision for
# execution in AWS CodeBuild when triggered by AWS CodePipeline and produces a
# source archive.
# The script assumes that you have the following environment variables set:
# * CODEBUILD_RESOLVED_SOURCE_VERSION - Set automatically when AWS CodePipeline
#   invokes AWS CodeBuild
# * GIT_REMOTE - A valid git remote that will be cloned
# This script is meant to be triggered only from branch-based executions of the
# pipeline.  It does not handle pull requests.

set -ex

[[ -d .git ]] && exit 1
[[ -z "${CODEBUILD_RESOLVED_SOURCE_VERSION}" ]] && exit 1
[[ -z "${GIT_REMOTE}" ]] && exit 1


mkdir archive
git clone "${GIT_REMOTE}" archive
cd archive
git checkout "${CODEBUILD_RESOLVED_SOURCE_VERSION}"
make release-tarball
make release-tarball-no-vendor
