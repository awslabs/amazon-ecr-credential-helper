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

# This script is an incredibly simple parser for the CHANGELOG.md in this
# repository.  It expects the changelog to have the following format:
#
# * Sections delimited by "#"
# * Section titles matching the VERSION file
#
# For example, consider the following sample changelog:
#
# # 1.2.3
# * Foo
# # 0.1.2-alpha
# * Bar
#
# If the VERSION file is set to "0.1.2-alpha", the output of this script is:
#
# # 0.1.2-alpha
# * Bar
set -e

# Normalize to working directory being source root (up two levels)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )
cd "${ROOT}"

VERSION="$(cat VERSION)"


match=""
while IFS= read -r line; do
    if [[ "${line:0:1}" == "#" ]]; then
        if [[ "${line}" == "# ${VERSION}" ]]; then
            match="y"
            continue
        else
            match=""
        fi
    fi
    if [[ -n "${match}" ]]; then
        echo "$line"
    fi
done < "CHANGELOG.md"
