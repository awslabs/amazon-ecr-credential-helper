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

# A script to generate a THIRD-PARTY-LICENSES file containing all the licenses that we use from third parties.
# NOTE: This only adds licenses from go dependencies. For other licenses, see NOTICE.

set -euo pipefail

# Normalize to working directory being root (up one level from ./scripts)
root=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

if test -d "${root}/ecr-login/vendor"; then
    echo "[ERROR]: generating THIRD-PARTY-LICENSES file while dependencies are vendored will result in unknown dependency versions in the licenses file."
    echo "[INFO]: To resolve, remove the vendored dependencies before generating the licenses file:"
    echo "----------------------------------------"
    echo "pushd ${root}"
    echo "rm -rf ecr-login/vendor"
    echo "make licenses"
    echo "git restore ecr-login/vendor"
    echo "popd"
    echo "----------------------------------------"
    exit 1
fi

license_file="${root}/THIRD-PARTY-LICENSES"

pushd "${root}/ecr-login"

# Remove content from the license file
truncate -s 0 "${license_file}"
{
    # The apache 2.0 license doesn't get modified with a copyright. To reduce duplication, add attribution for each project using the license, but include the license text just once.
    go-licenses report \
        --include_tests \
        --ignore github.com/awslabs/amazon-ecr-credential-helper \
        --template="${root}/scripts/third_party_licenses/apache.tpl" ./...
    cat "${root}/scripts/third_party_licenses/APACHE_LICENSE"
    # For other licenses, just use the entire license text from the package.
    go-licenses report \
        --include_tests \
        --ignore github.com/awslabs/amazon-ecr-credential-helper \
        --template="${root}/scripts/third_party_licenses/other.tpl" ./...
} >> "${license_file}"

popd
