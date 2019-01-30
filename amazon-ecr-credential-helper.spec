# Copyright 2018-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
Name:           amazon-ecr-credential-helper
Version:        0.3.0
Release:        1%{?dist}
Group:          Development/Tools
Vendor:         Amazon.com
License:        Apache 2.0
Summary:        Amazon ECR Docker Credential Helper
BuildArch:      x86_64 aarch64
BuildRoot:      ${_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

Source0: release.tar.gz

BuildRequires: golang >= 1.9

# The following 'Provides' lists the vendored dependencies bundled in
# and used to produce the amazon-ecr-credential-helper package. As dependencies
# are added or removed, this list should also be updated accordingly.
#
# You can use this to generate a list of the appropriate Provides
# statements by reading out the vendor directory:
#
# find ecr-login/vendor -name \*.go -exec dirname {} \; | sort | uniq | sed 's,^.*ecr-login/vendor/,,; s/^/bundled(golang(/; s/$/))/;' | sed 's/^/Provides:\t/' | expand -
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/awserr))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/awsutil))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/client))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/client/metadata))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/corehandlers))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/credentials))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/credentials/endpointcreds))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/credentials/processcreds))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/credentials/stscreds))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/csm))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/defaults))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/ec2metadata))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/endpoints))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/request))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/session))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/aws/signer/v4))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/internal/ini))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/internal/sdkio))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/internal/sdkrand))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/internal/sdkuri))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/internal/shareddefaults))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/json/jsonutil))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/jsonrpc))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/query))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/query/queryutil))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/rest))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/service/ecr))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/service/ecr/ecriface))
Provides:       bundled(golang(github.com/aws/aws-sdk-go/service/sts))
Provides:       bundled(golang(github.com/davecgh/go-spew/spew))
Provides:       bundled(golang(github.com/docker/docker-credential-helpers/credentials))
Provides:       bundled(golang(github.com/golang/mock/gomock))
Provides:       bundled(golang(github.com/jmespath/go-jmespath))
Provides:       bundled(golang(github.com/konsorten/go-windows-terminal-sequences))
Provides:       bundled(golang(github.com/mitchellh/go-homedir))
Provides:       bundled(golang(github.com/pkg/errors))
Provides:       bundled(golang(github.com/pmezard/go-difflib/difflib))
Provides:       bundled(golang(github.com/sirupsen/logrus))
Provides:       bundled(golang(github.com/stretchr/testify/assert))
Provides:       bundled(golang(golang.org/x/crypto/ssh/terminal))
Provides:       bundled(golang(golang.org/x/net/context))
Provides:       bundled(golang(golang.org/x/sys/unix))
Provides:       bundled(golang(golang.org/x/sys/windows))

%description
The Amazon ECR Docker Credential Helper is a credential helper for the Docker
daemon that makes it easier to use Amazon Elastic Container Registry.

%prep
%setup -c

%build
export GOPATH="$(pwd)/_gopath"
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
* Tue Jan 29 2019 Samuel Karp <skarp@amazon.com> - 0.3.0-1
- Added support for PrivateLink endpoints
* Tue Dec 4 2018 Samuel Karp <skarp@amazon.com> - 0.2.0-2
- Add aarch64 support
* Fri Nov 16 2018 Samuel Karp <skarp@amazon.com> - 0.2.0-1
- Initial packaging
