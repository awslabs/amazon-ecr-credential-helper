#!/usr/bin/env bash

#   Copyright The containerd Authors.

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

# A script to generate a THIRD_PARTY_LICENSES file containing all the licenses that we use from third parties.
# NOTE: This only adds licenses from go dependencies. For other licenses, see NOTICE.md

set -eux -o pipefail

CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOCI_SNAPSHOTTER_PROJECT_ROOT="${CUR_DIR}/.."
LICENSE_FILE="${SOCI_SNAPSHOTTER_PROJECT_ROOT}/THIRD_PARTY_LICENSES"

# Remove content from the license file
truncate -s 0 "${LICENSE_FILE}"
{
    # The apache 2.0 license doesn't get modified with a copywrite. To reduce duplication, add attribution for each project using the license, but include the license text just once.
    go-licenses report --template="${SOCI_SNAPSHOTTER_PROJECT_ROOT}/scripts/third_party_licenses/apache.tpl" --ignore github.com/awslabs/soci "${SOCI_SNAPSHOTTER_PROJECT_ROOT}"/...
    cat "${SOCI_SNAPSHOTTER_PROJECT_ROOT}/scripts/third_party_licenses/APACHE_LICENSE"
    # For other licenses, just use the entire license text from the package.
    go-licenses report --template="${SOCI_SNAPSHOTTER_PROJECT_ROOT}/scripts/third_party_licenses/other.tpl" --ignore github.com/awslabs/soci "${SOCI_SNAPSHOTTER_PROJECT_ROOT}"/...
} >> "${LICENSE_FILE}"
