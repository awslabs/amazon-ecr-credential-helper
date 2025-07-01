# Copyright 2018-2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the
# "License"). You may not use this file except in compliance
# with the License. A copy of the License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and
# limitations under the License.
%if 0%{?amzn} > 2
%define debug_package %{nil}
%endif
Name:           amazon-ecr-credential-helper
Version:        0.10.1
Release:        1%{?dist}
Group:          Development/Tools
Vendor:         Amazon.com
License:        Apache 2.0
Summary:        Amazon ECR Docker Credential Helper
BuildArch:      x86_64 aarch64
BuildRoot:      ${_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

Source0: release.tar.gz

BuildRequires: golang >= 1.23.10

# The following 'Provides' lists the vendored dependencies bundled in
# and used to produce the amazon-ecr-credential-helper package. As dependencies
# are added or removed, this list should also be updated accordingly.
#
# You can use this to generate a list of the appropriate Provides
# statements by reading out the vendor directory:
#
# find ecr-login/vendor -name \*.go -exec dirname {} \; | sort | uniq | sed 's,^.*ecr-login/vendor/,,; s/^/bundled(golang(/; s/$/))/;' | sed 's/^/Provides:\t/' | expand -
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/middleware))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/protocol/query))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/protocol/restjson))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/protocol/xml))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/ratelimit))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/retry))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/signer/v4))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/aws/transport/http))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/config))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/endpointcreds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/endpointcreds/internal/client))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/processcreds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/ssocreds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/credentials/stscreds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/feature/ec2/imds))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/config))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/endpoints))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/ini))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/rand))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/sdk))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/sdkio))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/strings))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/sync/singleflight))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/internal/timeconv))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecr))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecr/internal/endpoints))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecrpublic))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecrpublic/internal/endpoints))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecrpublic/types))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/ecr/types))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/internal/presigned-url))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sso))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sso/internal/endpoints))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sso/types))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sts))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sts/internal/endpoints))
Provides:     bundled(golang(github.com/aws/aws-sdk-go-v2/service/sts/types))
Provides:     bundled(golang(github.com/aws/smithy-go))
Provides:     bundled(golang(github.com/aws/smithy-go/encoding))
Provides:     bundled(golang(github.com/aws/smithy-go/encoding/httpbinding))
Provides:     bundled(golang(github.com/aws/smithy-go/encoding/json))
Provides:     bundled(golang(github.com/aws/smithy-go/encoding/xml))
Provides:     bundled(golang(github.com/aws/smithy-go/io))
Provides:     bundled(golang(github.com/aws/smithy-go/logging))
Provides:     bundled(golang(github.com/aws/smithy-go/middleware))
Provides:     bundled(golang(github.com/aws/smithy-go/ptr))
Provides:     bundled(golang(github.com/aws/smithy-go/rand))
Provides:     bundled(golang(github.com/aws/smithy-go/time))
Provides:     bundled(golang(github.com/aws/smithy-go/transport/http))
Provides:     bundled(golang(github.com/aws/smithy-go/transport/http/internal/io))
Provides:     bundled(golang(github.com/aws/smithy-go/waiter))
Provides:     bundled(golang(github.com/davecgh/go-spew/spew))
Provides:     bundled(golang(github.com/docker/docker-credential-helpers/credentials))
Provides:     bundled(golang(github.com/jmespath/go-jmespath))
Provides:     bundled(golang(github.com/konsorten/go-windows-terminal-sequences))
Provides:     bundled(golang(github.com/mitchellh/go-homedir))
Provides:     bundled(golang(github.com/pkg/errors))
Provides:     bundled(golang(github.com/pmezard/go-difflib/difflib))
Provides:     bundled(golang(github.com/sirupsen/logrus))
Provides:     bundled(golang(github.com/stretchr/testify/assert))
Provides:     bundled(golang(golang.org/x/sys/internal/unsafeheader))
Provides:     bundled(golang(golang.org/x/sys/unix))


%description
The Amazon ECR Docker Credential Helper is a credential helper for the Docker
daemon that makes it easier to use Amazon Elastic Container Registry.

%prep
%setup -c

%build
export GOPATH="$(pwd)/_gopath"
export GO111MODULE=off
mkdir -p "_gopath/src/github.com/awslabs"
ln -sv "$(pwd)" "_gopath/src/github.com/awslabs/amazon-ecr-credential-helper"
cd "_gopath/src/github.com/awslabs/amazon-ecr-credential-helper"
make

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/%{_bindir}
install -p -m 0755 \
  bin/local/docker-credential-ecr-login \
  %{buildroot}%{_bindir}/docker-credential-ecr-login
install -D -m 0644 \
  docs/docker-credential-ecr-login.1 \
  %{buildroot}%{_mandir}/man1/docker-credential-ecr-login.1

%files
%defattr(-,root,root,-)
%{_bindir}/docker-credential-ecr-login
%{_mandir}/man1/docker-credential-ecr-login.1*

%clean
rm -rf %{buildroot}

%changelog
* Tue Jul 01 2025 arjunry <arjunry@amazon.com> - 0.10.1-1
- Update to v0.10.1
'- Update spec to v0.10.1'

* Tue Jul 1 2025 Arjun Raja Yogidas <arjunry@amazon.com> - 0.10.1-1
- Update to v0.10.1
- fix CVE-2025-0913 and CVE-2025-4673
- Upgraded dependencies
* Wed Jun 4 2025 Shubhranshu Mahapatra <shubhum@amazon.com> - 0.10.0-1
- Update to v0.10.0
- Enhancement - Updated ECR pattern for ECR dual-stack endpoints for IPv6 support. 
- Upgraded dependencies
* Wed Sep 18 2024 Christopher R. Miller <milrchr@amazon.com> - 0.9.0-1
- Update to v0.9.0
- Enhancement - Updated ECR pattern to match C2S environments
- Enhancement - Added support for environment variable AWS_ECR_IGNORE_CREDS_STORAGE=true to ignore ADD and DELETE requests. This makes tools that try to docker login work with registries managed the amazon-ecr-credential-helper
- Enhancement - Updated ECR pattern for new isolated regions
- Upgraded dependencies
* Tue Aug 15 2023 Swagat Bora <sbora@amazon.com> - 0.7.1-1
- Allow callers to set log output
- Upgrade dependencies for bug fixes
* Fri Jan 14 2022 Austin Vazquez <macedonv@amazon.com> - 0.6.0-1
- Added support for AWS SSO
- Added support to assume roles via EC2 instance metadata
- Shared config file (~/.aws/config) is now always enabled. (AWS_SDK_LOAD_CONFIG environment variable is no longer supported)
* Mon Feb 15 2021 Samuel Karp <skarp@amazon.com> - 0.5.0-1
- Added support for ECR Public
- Added support for EC2 IMDSv2
- Enabled shared config file (~/.aws/config) by default
- Fixed bug with long credential_process responses
* Tue Jan 7 2020 Samuel Karp <skarp@amazon.com> - 0.4.0-1
- Added support for chaining assumed roles in the shared config file
- Added support for Web Identities and IAM Roles for Service Accounts (IRSA)
  with Kubernetes
- Log directory is now automatically created when the helper runs
* Tue Jan 29 2019 Samuel Karp <skarp@amazon.com> - 0.3.0-1
- Added support for PrivateLink endpoints
* Tue Dec 4 2018 Samuel Karp <skarp@amazon.com> - 0.2.0-2
- Add aarch64 support
* Fri Nov 16 2018 Samuel Karp <skarp@amazon.com> - 0.2.0-1
- Initial packaging
